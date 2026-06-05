package repository

import (
	"context"

	"github.com/google/uuid"
	"booking-service/internal/model"
)

func (r *Repo) CreateReview(ctx context.Context, review model.Review) error {
	_, err := r.db.Exec(ctx, `
		insert into public.reviews
		  (booking_id, client_id, master_id, rating, comment)
		values ($1,$2,$3,$4,$5)
		on conflict (booking_id) do update
		  set rating = excluded.rating,
		      comment = excluded.comment
	`, review.BookingID, review.ClientID, review.MasterID,
		review.Rating, review.Comment)
	return err
}

func (r *Repo) GetMasterReviews(ctx context.Context, masterID uuid.UUID) ([]model.Review, error) {
	rows, err := r.db.Query(ctx, `
		select id, booking_id, client_id, master_id, rating,
		       coalesce(comment,''), created_at
		from public.reviews
		where master_id = $1 and is_visible = true
		order by created_at desc
	`, masterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []model.Review
	for rows.Next() {
		var rv model.Review
		if err := rows.Scan(
			&rv.ID, &rv.BookingID, &rv.ClientID, &rv.MasterID,
			&rv.Rating, &rv.Comment, &rv.CreatedAt,
		); err != nil { return nil, err }
		list = append(list, rv)
	}
	return list, rows.Err()
}

func (r *Repo) HasReview(ctx context.Context, bookingID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`select count(*) from public.reviews where booking_id = $1`,
		bookingID).Scan(&count)
	return count > 0, err
}

func (r *Repo) CreateMasterReschedule(
	ctx context.Context,
	originalBookingID uuid.UUID,
	newBookingID      uuid.UUID,
	clientID          uuid.UUID,
	masterID          uuid.UUID,
	reason            string,
) error {
	_, err := r.db.Exec(ctx, `
		insert into public.master_reschedules
		  (original_booking_id, new_booking_id, client_id, master_id, reason)
		values ($1,$2,$3,$4,$5)
	`, originalBookingID, newBookingID, clientID, masterID, reason)
	return err
}

func (r *Repo) GetClientReschedules(
	ctx context.Context,
	clientID uuid.UUID,
) ([]model.MasterReschedule, error) {
	rows, err := r.db.Query(ctx, `
		select
			mr.id, mr.original_booking_id, mr.new_booking_id,
			mr.client_id, mr.master_id,
			coalesce(mr.reason,''), mr.status, mr.created_at,
			coalesce(s.name,''),
			coalesce(p.full_name,''),
			nb.starts_at,
			nb.price_paid
		from public.master_reschedules mr
		join public.bookings nb on nb.id = mr.new_booking_id
		join public.services  s on  s.id = nb.service_id
		join public.profiles  p on  p.id = nb.master_id
		where mr.client_id = $1
		  and mr.status = 'pending'
		order by mr.created_at desc
	`, clientID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []model.MasterReschedule
	for rows.Next() {
		var m model.MasterReschedule
		if err := rows.Scan(
			&m.ID, &m.OriginalBookingID, &m.NewBookingID,
			&m.ClientID, &m.MasterID,
			&m.Reason, &m.Status, &m.CreatedAt,
			&m.ServiceName, &m.NewMasterName,
			&m.NewStartsAt, &m.NewPrice,
		); err != nil { return nil, err }
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *Repo) DeclineMasterReschedule(
	ctx context.Context,
	id       string,
	clientID uuid.UUID,
) error {
	var newBookingID uuid.UUID
	err := r.db.QueryRow(ctx, `
		update public.master_reschedules
		set status = 'declined'
		where id = $1 and client_id = $2
		returning new_booking_id
	`, id, clientID).Scan(&newBookingID)
	if err != nil { return err }

	_, err = r.db.Exec(ctx,
		`update public.bookings set status = 'cancelled' where id = $1`,
		newBookingID)
	return err
}

func (r *Repo) DeleteReviewByMaster(ctx context.Context, reviewID string, masterID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		update public.reviews
		set is_visible = false
		where id = $1::uuid and master_id = $2
	`, reviewID, masterID)
	return err
}

func (r *Repo) DeleteReviewByAdmin(ctx context.Context, reviewID string) error {
	_, err := r.db.Exec(ctx, `
		update public.reviews
		set is_visible = false
		where id = $1::uuid
	`, reviewID)
	return err
}