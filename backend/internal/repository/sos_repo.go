package repository

import (
	"context"

	"github.com/google/uuid"
	"booking-service/internal/model"
)

func (r *Repo) CreateSOSRequest(ctx context.Context, req model.SOSRequest) (*model.SOSRequest, error) {
	req.ID = uuid.New()
	err := r.db.QueryRow(ctx, `
		insert into public.sos_requests
		  (id, client_id, master_id, service_id,
		   preferred_range_start, preferred_range_end,
		   base_price, sos_price, client_note)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		returning id, created_at, expires_at
	`, req.ID, req.ClientID, req.MasterID, req.ServiceID,
		req.PreferredRangeStart, req.PreferredRangeEnd,
		req.BasePrice, req.SOSPrice, req.ClientNote,
	).Scan(&req.ID, &req.CreatedAt, &req.ExpiresAt)
	return &req, err
}

func (r *Repo) GetClientSOSRequests(ctx context.Context, clientID uuid.UUID) ([]model.SOSRequest, error) {
	rows, err := r.db.Query(ctx, `
		select s.id, s.client_id, s.master_id, s.service_id,
		       s.preferred_range_start, s.preferred_range_end,
		       s.base_price, s.sos_price, s.status,
		       coalesce(s.client_note,''), coalesce(s.master_note,''),
		       s.expires_at, s.created_at
		from public.sos_requests s
		where s.client_id = $1
		order by s.created_at desc limit 20
	`, clientID)
	if err != nil { return nil, err }
	defer rows.Close()
	return scanSOSRows(rows)
}

func (r *Repo) GetMasterSOSRequests(ctx context.Context, masterID uuid.UUID) ([]model.SOSRequest, error) {
	rows, err := r.db.Query(ctx, `
		select s.id, s.client_id, s.master_id, s.service_id,
		       s.preferred_range_start, s.preferred_range_end,
		       s.base_price, s.sos_price, s.status,
		       coalesce(s.client_note,''), coalesce(s.master_note,''),
		       s.expires_at, s.created_at,
		       p.full_name as client_name,
		       svc.name as service_name
		from public.sos_requests s
		join public.profiles p on p.id = s.client_id
		join public.services svc on svc.id = s.service_id
		where s.master_id = $1 and s.status = 'pending'
		  and s.expires_at > now()
		order by s.created_at desc
	`, masterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []model.SOSRequest
	for rows.Next() {
		var s model.SOSRequest
		if err := rows.Scan(
			&s.ID, &s.ClientID, &s.MasterID, &s.ServiceID,
			&s.PreferredRangeStart, &s.PreferredRangeEnd,
			&s.BasePrice, &s.SOSPrice, &s.Status,
			&s.ClientNote, &s.MasterNote,
			&s.ExpiresAt, &s.CreatedAt,
			&s.ClientName, &s.ServiceName,
		); err != nil { return nil, err }
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *Repo) RespondToSOS(
    ctx context.Context,
    sosID, masterID uuid.UUID,
    accept bool,
    note string,
) error {
    status := "declined"
    if accept { status = "accepted" }

    _, err := r.db.Exec(ctx, `
        update public.sos_requests
        set status = $1, master_note = $2
        where id = $3 and master_id = $4
    `, status, note, sosID, masterID)
    return err
}

func (r *Repo) GetSOSRequest(ctx context.Context, id uuid.UUID) (*model.SOSRequest, error) {
	var s model.SOSRequest
	err := r.db.QueryRow(ctx, `
		select id, client_id, master_id, service_id,
		       preferred_range_start, preferred_range_end,
		       base_price, sos_price, status,
		       coalesce(client_note,''), coalesce(master_note,''),
		       expires_at, created_at
		from public.sos_requests where id = $1
	`, id).Scan(
		&s.ID, &s.ClientID, &s.MasterID, &s.ServiceID,
		&s.PreferredRangeStart, &s.PreferredRangeEnd,
		&s.BasePrice, &s.SOSPrice, &s.Status,
		&s.ClientNote, &s.MasterNote,
		&s.ExpiresAt, &s.CreatedAt,
	)
	return &s, err
}

func scanSOSRows(rows interface{ Next() bool; Scan(...any) error; Err() error }) ([]model.SOSRequest, error) {
	var list []model.SOSRequest
	for rows.Next() {
		var s model.SOSRequest
		if err := rows.Scan(
			&s.ID, &s.ClientID, &s.MasterID, &s.ServiceID,
			&s.PreferredRangeStart, &s.PreferredRangeEnd,
			&s.BasePrice, &s.SOSPrice, &s.Status,
			&s.ClientNote, &s.MasterNote,
			&s.ExpiresAt, &s.CreatedAt,
		); err != nil { return nil, err }
		list = append(list, s)
	}
	return list, rows.Err()
}