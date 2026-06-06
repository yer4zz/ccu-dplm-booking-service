package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "booking-service/internal/model"
)

func (r *Repo) IsSlotAvailable(ctx context.Context, masterID uuid.UUID, start, end time.Time) (bool, error) {
    var count int
    err := r.db.QueryRow(ctx, `
        SELECT COUNT(*) FROM public.bookings
        WHERE master_id = $1
          AND status NOT IN ('cancelled', 'no_show')
          AND starts_at < $3
          AND ends_at   > $2
    `, masterID, start, end).Scan(&count)
    if err != nil {
        return false, err
    }
    return count == 0, nil
}

func (r *Repo) CreateBooking(ctx context.Context, b model.Booking) (*model.Booking, error) {
    b.ID = uuid.New()
    err := r.db.QueryRow(ctx, `
        INSERT INTO public.bookings
            (id, client_id, master_id, service_id, starts_at, ends_at, status, price_paid, notes)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        RETURNING id, created_at
    `, b.ID, b.ClientID, b.MasterID, b.ServiceID,
       b.StartsAt, b.EndsAt, b.Status, b.PricePaid, b.Notes,
    ).Scan(&b.ID, &b.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &b, nil
}

func (r *Repo) GetBooking(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
    var b model.Booking
    err := r.db.QueryRow(ctx, `
        SELECT id, client_id, master_id, service_id,
               starts_at, ends_at, status, price_paid, notes, created_at
        FROM public.bookings WHERE id = $1
    `, id).Scan(
        &b.ID, &b.ClientID, &b.MasterID, &b.ServiceID,
        &b.StartsAt, &b.EndsAt, &b.Status, &b.PricePaid, &b.Notes, &b.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &b, nil
}

func (r *Repo) GetBookingsForDay(ctx context.Context, masterID uuid.UUID, date time.Time) ([]model.Booking, error) {
    dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    dayEnd   := dayStart.Add(24 * time.Hour)

    rows, err := r.db.Query(ctx, `
        SELECT id, client_id, master_id, service_id,
               starts_at, ends_at, status, price_paid, notes, created_at
        FROM public.bookings
        WHERE master_id = $1
          AND starts_at >= $2 AND starts_at < $3
          AND status NOT IN ('cancelled', 'no_show')
        ORDER BY starts_at
    `, masterID, dayStart, dayEnd)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.Booking
    for rows.Next() {
        var b model.Booking
        if err := rows.Scan(
            &b.ID, &b.ClientID, &b.MasterID, &b.ServiceID,
            &b.StartsAt, &b.EndsAt, &b.Status, &b.PricePaid, &b.Notes, &b.CreatedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, b)
    }
    return list, rows.Err()
}

func (r *Repo) GetMasterBookings(ctx context.Context, masterID uuid.UUID) ([]model.Booking, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.client_id, b.master_id, b.service_id,
		       b.starts_at, b.ends_at, b.status, b.price_paid,
		       COALESCE(b.notes, ''), COALESCE(b.vibe_mode, ''), b.created_at,
		       COALESCE(s.name, ''), COALESCE(p.full_name, '')
		FROM public.bookings b
		LEFT JOIN public.services s ON s.id = b.service_id
		LEFT JOIN public.profiles p ON p.id = b.client_id
		WHERE b.master_id = $1
		ORDER BY b.starts_at DESC
	`, masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Booking
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(
			&b.ID, &b.ClientID, &b.MasterID, &b.ServiceID,
			&b.StartsAt, &b.EndsAt, &b.Status, &b.PricePaid,
			&b.Notes, &b.VibeMode, &b.CreatedAt,
			&b.ServiceName, &b.ClientName,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (r *Repo) GetClientBookings(ctx context.Context, clientID uuid.UUID) ([]model.Booking, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.client_id, b.master_id, b.service_id,
		       b.starts_at, b.ends_at, b.status, b.price_paid,
		       COALESCE(b.notes, ''), b.created_at,
		       COALESCE(s.name, ''), COALESCE(p.full_name, '')
		FROM public.bookings b
		LEFT JOIN public.services s ON s.id = b.service_id
		LEFT JOIN public.profiles p ON p.id = b.master_id
		WHERE b.client_id = $1
		ORDER BY b.starts_at DESC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Booking
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(
			&b.ID, &b.ClientID, &b.MasterID, &b.ServiceID,
			&b.StartsAt, &b.EndsAt, &b.Status, &b.PricePaid,
			&b.Notes, &b.CreatedAt,
			&b.ServiceName, &b.MasterName,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (r *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.BookingStatus) error {
    _, err := r.db.Exec(ctx, `
        UPDATE public.bookings SET status = $1 WHERE id = $2
    `, status, id)
    return err
}

func (r *Repo) GetBookingNotification(ctx context.Context, bookingID uuid.UUID) (*model.BookingNotification, error) {
	var n model.BookingNotification
	err := r.db.QueryRow(ctx, `
		SELECT
			b.id::text,
			COALESCE(cp.full_name, ''),
			COALESCE(cu.email, ''),
			COALESCE(mp.full_name, ''),
			COALESCE(mu.email, ''),
			COALESCE(s.name, ''),
			b.starts_at,
			COALESCE(s.duration_min, 0),
			COALESCE(b.price_paid, 0)
		FROM public.bookings b
		LEFT JOIN public.profiles cp ON cp.id = b.client_id
		LEFT JOIN auth.users    cu ON cu.id = b.client_id
		LEFT JOIN public.profiles mp ON mp.id = b.master_id
		LEFT JOIN auth.users    mu ON mu.id = b.master_id
		LEFT JOIN public.services s ON s.id = b.service_id
		WHERE b.id = $1
	`, bookingID).Scan(
		&n.BookingID,
		&n.ClientName, &n.ClientEmail,
		&n.MasterName, &n.MasterEmail,
		&n.ServiceName,
		&n.StartsAt, &n.DurationMin, &n.Price,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}


func (r *Repo) AddBookingService(ctx context.Context, bookingID, serviceID uuid.UUID, price float64) error {
	_, err := r.db.Exec(ctx, `
		insert into public.booking_services (booking_id, service_id, price)
		values ($1, $2, $3) on conflict do nothing
	`, bookingID, serviceID, price)
	return err
}

func (r *Repo) CreateRescheduleOffer(ctx context.Context, bookingID, clientID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		insert into public.reschedule_offers (original_booking_id, client_id)
		values ($1, $2)
	`, bookingID, clientID)
	return err
}

func (r *Repo) GetRescheduleOffers(ctx context.Context, clientID uuid.UUID) ([]model.RescheduleOffer, error) {
	rows, err := r.db.Query(ctx, `
		select id, original_booking_id, client_id, status, offered_at, expires_at
		from public.reschedule_offers
		where client_id = $1 and status = 'pending'
		  and expires_at > now()
		order by offered_at desc
	`, clientID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []model.RescheduleOffer
	for rows.Next() {
		var o model.RescheduleOffer
		if err := rows.Scan(
			&o.ID, &o.OriginalBookingID, &o.ClientID,
			&o.Status, &o.OfferedAt, &o.ExpiresAt,
		); err != nil { return nil, err }
		list = append(list, o)
	}
	return list, rows.Err()
}

func (r *Repo) RescheduleBooking(
	ctx context.Context,
	bookingID uuid.UUID,
	newMasterID *uuid.UUID,
	newStartsAt time.Time,
	newEndsAt   time.Time,
	newPrice    float64,
) error {
	masterID := bookingID
	_ = masterID

	query := `
		update public.bookings
		set starts_at  = $1,
		    ends_at    = $2,
		    status     = 'pending',
		    price_paid = $3
	`
	args := []any{newStartsAt, newEndsAt, newPrice}

	if newMasterID != nil {
		query += `, master_id = $4 where id = $5`
		args = append(args, *newMasterID, bookingID)
	} else {
		query += ` where id = $4`
		args = append(args, bookingID)
	}

	_, err := r.db.Exec(ctx, query, args...)
	return err
}



func (r *Repo) GetBookingServiceIDs(ctx context.Context, bookingID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
		select service_id from public.booking_services
		where booking_id = $1
	`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *Repo) IsSlotAvailableExclude(
	ctx context.Context,
	masterID uuid.UUID,
	start, end time.Time,
	excludeBookingID uuid.UUID,
) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		select count(*) from public.bookings
		where master_id = $1
		  and id != $4
		  and status not in ('cancelled', 'no_show')
		  and starts_at < $3
		  and ends_at   > $2
	`, masterID, start, end, excludeBookingID).Scan(&count)
	return count == 0, err
}

func (r *Repo) CreateAutoReschedule(
	ctx context.Context,
	originalBookingID uuid.UUID,
	newBookingID      uuid.UUID,
	clientID          uuid.UUID,
) error {
	_, err := r.db.Exec(ctx, `
		insert into public.auto_reschedules
		  (original_booking_id, new_booking_id, client_id)
		values ($1, $2, $3)
	`, originalBookingID, newBookingID, clientID)
	return err
}


func (r *Repo) DeclineAutoReschedule(ctx context.Context, id string, clientID uuid.UUID) error {
	var newBookingID uuid.UUID
	err := r.db.QueryRow(ctx, `
		update public.auto_reschedules
		set status = 'declined'
		where id = $1 and client_id = $2
		returning new_booking_id
	`, id, clientID).Scan(&newBookingID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `
		update public.bookings set status = 'cancelled'
		where id = $1
	`, newBookingID)
	return err
}