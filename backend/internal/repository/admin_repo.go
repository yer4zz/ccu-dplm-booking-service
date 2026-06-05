package repository

import (
	"context"

	"github.com/google/uuid"
	"booking-service/internal/model"
)

func (r *Repo) GetAllBookings(ctx context.Context) ([]model.Booking, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, client_id, master_id, service_id,
		       starts_at, ends_at, status, price_paid, notes, created_at
		FROM public.bookings
		ORDER BY starts_at DESC
		LIMIT 500
	`)
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

func (r *Repo) SetServiceActive(ctx context.Context, id uuid.UUID, active bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE public.services SET is_active = $1 WHERE id = $2`, active, id)
	return err
}

func (r *Repo) SetMasterActive(ctx context.Context, id uuid.UUID, active bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE public.masters SET is_active = $1 WHERE id = $2`, active, id)
	return err
}