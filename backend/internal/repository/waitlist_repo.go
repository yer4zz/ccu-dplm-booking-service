package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"booking-service/internal/model"
)

func (r *Repo) AddToWaitlist(ctx context.Context, e model.WaitlistEntry) (*model.WaitlistEntry, error) {
	e.ID = uuid.New()
	err := r.db.QueryRow(ctx, `
		insert into public.waitlist
		  (id, client_id, master_id, service_id, preferred_date)
		values ($1,$2,$3,$4,$5)
		returning id, created_at
	`, e.ID, e.ClientID, e.MasterID, e.ServiceID, e.PreferredDate,
	).Scan(&e.ID, &e.CreatedAt)
	return &e, err
}

func (r *Repo) GetWaitlist(ctx context.Context, clientID uuid.UUID) ([]model.WaitlistEntry, error) {
	rows, err := r.db.Query(ctx, `
		select id, client_id, master_id, service_id, preferred_date, status, created_at
		from public.waitlist
		where client_id = $1 and status = 'waiting'
		order by created_at desc
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.WaitlistEntry
	for rows.Next() {
		var e model.WaitlistEntry
		if err := rows.Scan(
			&e.ID, &e.ClientID, &e.MasterID, &e.ServiceID,
			&e.PreferredDate, &e.Status, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *Repo) NotifyWaitlistForSlot(ctx context.Context, masterID, serviceID uuid.UUID, slotTime time.Time) (*model.WaitlistEntry, error) {
	var e model.WaitlistEntry
	err := r.db.QueryRow(ctx, `
		update public.waitlist
		set status = 'notified',
		    notified_at = now(),
		    expires_at = now() + interval '15 minutes'
		where id = (
		  select id from public.waitlist
		  where master_id = $1
		    and service_id = $2
		    and status = 'waiting'
		  order by created_at asc
		  limit 1
		)
		returning id, client_id, master_id, service_id, status
	`, masterID, serviceID).Scan(
		&e.ID, &e.ClientID, &e.MasterID, &e.ServiceID, &e.Status,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repo) RemoveFromWaitlist(ctx context.Context, id, clientID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		update public.waitlist set status = 'expired'
		where id = $1 and client_id = $2
	`, id, clientID)
	return err
}