package repository

import (
    "context"

    "github.com/google/uuid"
    "booking-service/internal/model"
)

func (r *Repo) ListServices(ctx context.Context) ([]model.Service, error) {
    rows, err := r.db.Query(ctx, `
        SELECT id, name, category, description, duration_min, price, is_active, sort_order,
               coalesce(name_kz, name), coalesce(name_en, name),
               coalesce(desc_kz, ''),   coalesce(desc_en, '')
        FROM public.services
        WHERE is_active = true
        ORDER BY sort_order, name
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.Service
    for rows.Next() {
        var s model.Service
        if err := rows.Scan(
            &s.ID, &s.Name, &s.Category, &s.Description,
            &s.DurationMin, &s.Price, &s.IsActive, &s.SortOrder,
            &s.NameKz, &s.NameEn, &s.DescKz, &s.DescEn,
        ); err != nil {
            return nil, err
        }
        list = append(list, s)
    }
    return list, rows.Err()
}

func (r *Repo) GetMasterService(ctx context.Context, masterID, serviceID uuid.UUID) (*model.MasterService, error) {
    var ms model.MasterService
    err := r.db.QueryRow(ctx, `
        SELECT s.id, s.name, s.category, s.description, s.duration_min, s.price, s.is_active, s.sort_order,
               ms.custom_price
        FROM public.master_services ms
        JOIN public.services s ON s.id = ms.service_id
        WHERE ms.master_id = $1 AND ms.service_id = $2
    `, masterID, serviceID).Scan(
        &ms.ID, &ms.Name, &ms.Category, &ms.Description,
        &ms.DurationMin, &ms.Price, &ms.IsActive, &ms.SortOrder,
        &ms.CustomPrice,
    )
    if err != nil {
        return nil, err
    }
    return &ms, nil
}

func (r *Repo) GetServicePrice(ctx context.Context, masterID, serviceID uuid.UUID) (float64, int, error) {
	var price float64
	var duration int
	err := r.db.QueryRow(ctx, `
		select
			coalesce(ms.custom_price, s.price) as price,
			s.duration_min
		from public.services s
		left join public.master_services ms
			on ms.service_id = s.id and ms.master_id = $1
		where s.id = $2 and s.is_active = true
	`, masterID, serviceID).Scan(&price, &duration)
	return price, duration, err
}