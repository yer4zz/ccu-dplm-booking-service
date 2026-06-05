package repository

import (
	"context"

	"github.com/google/uuid"
	"booking-service/internal/model"
)

func (r *Repo) GetLoyaltyAccount(ctx context.Context, clientID uuid.UUID) (*model.LoyaltyAccount, error) {
	var a model.LoyaltyAccount
	err := r.db.QueryRow(ctx, `
		select client_id, balance, total_earned, updated_at
		from public.loyalty_accounts
		where client_id = $1
	`, clientID).Scan(&a.ClientID, &a.Balance, &a.TotalEarned, &a.UpdatedAt)
	if err != nil {
		_, err2 := r.db.Exec(ctx,
			`insert into public.loyalty_accounts (client_id) values ($1) on conflict do nothing`,
			clientID)
		if err2 != nil {
			return nil, err2
		}
		return &model.LoyaltyAccount{ClientID: clientID, Balance: 0}, nil
	}
	return &a, nil
}

func (r *Repo) GetPointTransactions(ctx context.Context, clientID uuid.UUID) ([]model.PointTransaction, error) {
	rows, err := r.db.Query(ctx, `
		select id, client_id, booking_id, type, amount, description, created_at
		from public.point_transactions
		where client_id = $1
		order by created_at desc
		limit 50
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PointTransaction
	for rows.Next() {
		var t model.PointTransaction
		if err := rows.Scan(
			&t.ID, &t.ClientID, &t.BookingID,
			&t.Type, &t.Amount, &t.Description, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *Repo) IsFirstBooking(ctx context.Context, clientID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		select count(*) from public.bookings
		where client_id = $1 and status not in ('cancelled')
	`, clientID).Scan(&count)
	return count == 0, err
}

func (r *Repo) DeductPoints(ctx context.Context, clientID, bookingID uuid.UUID, points int) error {
	_, err := r.db.Exec(ctx, `
		update public.loyalty_accounts
		set balance = balance - $2, updated_at = now()
		where client_id = $1 and balance >= $2
	`, clientID, points)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `
		insert into public.point_transactions
		  (client_id, booking_id, type, amount, description)
		values ($1, $2, 'redeem', $3, 'Списание баллов при бронировании')
	`, clientID, bookingID, points)
	return err
}

func (r *Repo) SaveDiscount(ctx context.Context, d model.Discount) error {
	_, err := r.db.Exec(ctx, `
		insert into public.discounts
		  (booking_id, type, percent, points_used, amount)
		values ($1, $2, $3, $4, $5)
		on conflict (booking_id) do update
		  set percent = excluded.percent,
		      points_used = excluded.points_used,
		      amount = excluded.amount
	`, d.BookingID, d.Type, d.Percent, d.PointsUsed, d.Amount)
	return err
}