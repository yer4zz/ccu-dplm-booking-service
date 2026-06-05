package repository

import (
	"booking-service/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repo) ListMasters(ctx context.Context) ([]model.Master, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			m.id::text,
			p.full_name,
			COALESCE(m.bio, ''),
			COALESCE(m.experience_years, 0),
			COALESCE(p.avatar_url, ''),
			COALESCE(m.instagram, ''),
			m.is_active,
			COALESCE(s.id::text, ''),
			COALESCE(s.name, ''),
			COALESCE(s.duration_min, 0),
			COALESCE(s.price, 0)
		FROM public.masters m
		JOIN public.profiles p ON p.id = m.id
		LEFT JOIN public.master_services ms ON ms.master_id = m.id
		LEFT JOIN public.services s ON s.id = ms.service_id AND s.is_active = true
		WHERE m.is_active = true
		ORDER BY p.full_name, s.name
	`)
	if err != nil {
		return nil, fmt.Errorf("ListMasters query: %w", err)
	}
	defer rows.Close()

	mastersMap := make(map[string]*model.Master)
	var order []string

	for rows.Next() {
		var (
			masterID   string
			fullName   string
			bio        string
			expYears   int
			avatarURL  string
			instagram  string
			isActive   bool
			serviceID  string
			svcName    string
			svcDur     int
			svcPrice   float64
		)
		if err := rows.Scan(
			&masterID, &fullName, &bio, &expYears,
			&avatarURL, &instagram, &isActive,
			&serviceID, &svcName, &svcDur, &svcPrice,
		); err != nil {
			return nil, fmt.Errorf("ListMasters scan: %w", err)
		}

		if _, exists := mastersMap[masterID]; !exists {
			uid, _ := uuid.Parse(masterID)
			mastersMap[masterID] = &model.Master{
				ID:              uid,
				FullName:        fullName,
				Bio:             bio,
				ExperienceYears: expYears,
				AvatarURL:       avatarURL,
				Instagram:       instagram,
				IsActive:        isActive,
				Services:        []model.Service{},
			}
			order = append(order, masterID)
		}

		if serviceID != "" {
			svcUID, _ := uuid.Parse(serviceID)
			mastersMap[masterID].Services = append(
				mastersMap[masterID].Services,
				model.Service{
					ID:          svcUID,
					Name:        svcName,
					DurationMin: svcDur,
					Price:       svcPrice,
				},
			)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListMasters rows: %w", err)
	}

	result := make([]model.Master, 0, len(order))
	for _, id := range order {
		result = append(result, *mastersMap[id])
	}
	return result, nil
}

func (r *Repo) getMasterServices(ctx context.Context, masterID string) ([]model.Service, error) {
    rows, err := r.db.Query(ctx, `
        SELECT s.id, s.name, s.category, s.description, s.duration_min,
               COALESCE(ms.custom_price, s.price) as price, s.is_active, s.sort_order
        FROM public.master_services ms
        JOIN public.services s ON s.id = ms.service_id
        WHERE ms.master_id = $1 AND s.is_active = true
    `, masterID)
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
        ); err != nil {
            return nil, err
        }
        list = append(list, s)
    }
    return list, rows.Err()
}

func (r *Repo) GetDaySchedule(ctx context.Context, masterID string, dayOfWeek int) (*model.Schedule, error) {
    var s model.Schedule
    err := r.db.QueryRow(ctx, `
        SELECT master_id, day_of_week,
               to_char(start_time, 'HH24:MI') as start_time,
               to_char(end_time,   'HH24:MI') as end_time
        FROM public.master_schedules
        WHERE master_id = $1 AND day_of_week = $2
    `, masterID, dayOfWeek).Scan(&s.MasterID, &s.DayOfWeek, &s.StartTime, &s.EndTime)
    if err != nil {
        return nil, err
    }
    return &s, nil
}