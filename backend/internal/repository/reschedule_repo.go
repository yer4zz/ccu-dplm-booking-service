package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AutoReschedule struct {
	ID                 string    `json:"id"`
	OriginalBookingID  string    `json:"original_booking_id"`
	NewBookingID       string    `json:"new_booking_id"`
	ClientID           string    `json:"client_id"`
	NewMasterName      string    `json:"new_master_name"`
	ServiceName        string    `json:"service_name"`
	NewStartsAt        time.Time `json:"new_starts_at"`
	PricePaid          float64   `json:"price_paid"`
	CreatedAt          time.Time `json:"created_at"`
}

func (r *Repo) FindAlternativeMaster(
	ctx context.Context,
	excludeMasterID uuid.UUID,
	serviceID       uuid.UUID,
	startsAt        time.Time,
	endsAt          time.Time,
) (*uuid.UUID, error) {
	var masterID uuid.UUID
	err := r.db.QueryRow(ctx, `
		select m.id
		from public.masters m
		where m.is_active = true
		  and m.id != $1
		  and exists (
		    select 1 from public.services s
		    where s.id = $2 and s.is_active = true
		  )
		  and not exists (
		    select 1 from public.bookings b
		    where b.master_id = m.id
		      and b.status not in ('cancelled','no_show')
		      and b.starts_at < $4
		      and b.ends_at   > $3
		  )
		  and exists (
		    select 1 from public.master_schedules sch
		    where sch.master_id   = m.id
		      and sch.day_of_week = extract(dow from $3::timestamptz)::int
		      and sch.start_time  <= $3::time
		      and sch.end_time    >= $4::time
		  )
		limit 1
	`, excludeMasterID, serviceID, startsAt, endsAt).Scan(&masterID)

	if err != nil {
		return nil, err
	}
	return &masterID, nil
}

func (r *Repo) SaveAutoReschedule(
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

func (r *Repo) GetAutoReschedules(
	ctx context.Context,
	clientID uuid.UUID,
) ([]AutoReschedule, error) {
	rows, err := r.db.Query(ctx, `
		select
			ar.id::text,
			ar.original_booking_id::text,
			ar.new_booking_id::text,
			ar.client_id::text,
			coalesce(p.full_name, ''),
			coalesce(s.name, ''),
			nb.starts_at,
			nb.price_paid,
			ar.created_at
		from public.auto_reschedules ar
		join public.bookings  nb on nb.id = ar.new_booking_id
		join public.profiles   p on  p.id = nb.master_id
		join public.services   s on  s.id = nb.service_id
		where ar.client_id = $1
		order by ar.created_at desc
		limit 20
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AutoReschedule
	for rows.Next() {
		var a AutoReschedule
		if err := rows.Scan(
			&a.ID, &a.OriginalBookingID, &a.NewBookingID, &a.ClientID,
			&a.NewMasterName, &a.ServiceName,
			&a.NewStartsAt, &a.PricePaid, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *Repo) CancelNewBookingFromReschedule(
	ctx context.Context,
	autoRescheduleID string,
	clientID         uuid.UUID,
) error {
	var newBookingID uuid.UUID
	err := r.db.QueryRow(ctx, `
		select new_booking_id from public.auto_reschedules
		where id = $1 and client_id = $2
	`, autoRescheduleID, clientID).Scan(&newBookingID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		update public.bookings set status = 'cancelled'
		where id = $1
	`, newBookingID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		delete from public.auto_reschedules
		where id = $1
	`, autoRescheduleID)
	return err
}