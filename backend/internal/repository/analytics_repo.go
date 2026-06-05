package repository

import (
	"context"
	"time"
)

type DayRevenue struct {
	Day           time.Time `db:"day"             json:"day"`
	Revenue       float64   `db:"revenue"         json:"revenue"`
	BookingsCount int64     `db:"bookings_count"  json:"bookings_count"`
}

type ServiceStat struct {
	ServiceName string  `db:"service_name" json:"service_name"`
	Count       int64   `db:"count"        json:"count"`
	Revenue     float64 `db:"revenue"      json:"revenue"`
}

type MasterStat struct {
	MasterName string  `db:"master_name" json:"master_name"`
	Total      int64   `db:"total"       json:"total"`
	Completed  int64   `db:"completed"   json:"completed"`
	Revenue    float64 `db:"revenue"     json:"revenue"`
	AvgPrice   float64 `db:"avg_price"   json:"avg_price"`
}

type HourlyLoad struct {
	HourOfDay int   `db:"hour_of_day" json:"hour"`
	DayOfWeek int   `db:"day_of_week" json:"dow"`
	Count     int64 `db:"count"       json:"count"`
}

type OverviewStats struct {
	Revenue    float64 `json:"revenue"`
	Total      int64   `json:"total"`
	Completed  int64   `json:"completed"`
	Cancelled  int64   `json:"cancelled"`
	NoShow     int64   `json:"no_show"`
	Pending    int64   `json:"pending"`
	AvgPrice   float64 `json:"avg_price"`
	CancelRate float64 `json:"cancel_rate"`
}

func (r *Repo) GetOverviewStats(ctx context.Context, from, to time.Time) (*OverviewStats, error) {
	var s OverviewStats
	err := r.db.QueryRow(ctx, `
		select
			coalesce(sum(price_paid) filter (where status='completed'),0),
			count(*),
			count(*) filter (where status='completed'),
			count(*) filter (where status='cancelled'),
			count(*) filter (where status='no_show'),
			count(*) filter (where status in ('pending','confirmed')),
			coalesce(avg(price_paid) filter (where status='completed'),0),
			coalesce(round(
				count(*) filter (where status='cancelled')::numeric /
				nullif(count(*),0)*100, 1
			),0)
		from public.bookings
		where starts_at >= $1 and starts_at < $2
	`, from, to).Scan(
		&s.Revenue, &s.Total, &s.Completed,
		&s.Cancelled, &s.NoShow, &s.Pending,
		&s.AvgPrice, &s.CancelRate,
	)
	return &s, err
}

func (r *Repo) GetRevenueByDay(ctx context.Context, from, to time.Time) ([]DayRevenue, error) {
	rows, err := r.db.Query(ctx, `
		select
			starts_at::date as day,
			coalesce(sum(price_paid) filter (where status='completed'),0) as revenue,
			count(*) as bookings_count
		from public.bookings
		where starts_at >= $1 and starts_at < $2
		group by starts_at::date
		order by day
	`, from, to)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []DayRevenue
	for rows.Next() {
		var d DayRevenue
		if err := rows.Scan(&d.Day, &d.Revenue, &d.BookingsCount); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *Repo) GetTopServices(ctx context.Context, from, to time.Time) ([]ServiceStat, error) {
	rows, err := r.db.Query(ctx, `
		select s.name,
			count(*) as count,
			coalesce(sum(b.price_paid) filter (where b.status='completed'),0) as revenue
		from public.bookings b
		join public.services s on s.id = b.service_id
		where b.starts_at >= $1 and b.starts_at < $2
		group by s.name
		order by count desc limit 7
	`, from, to)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []ServiceStat
	for rows.Next() {
		var s ServiceStat
		if err := rows.Scan(&s.ServiceName, &s.Count, &s.Revenue); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *Repo) GetMasterStats(ctx context.Context, from, to time.Time) ([]MasterStat, error) {
	rows, err := r.db.Query(ctx, `
		select p.full_name,
			count(*) as total,
			count(*) filter (where b.status='completed') as completed,
			coalesce(sum(b.price_paid) filter (where b.status='completed'),0) as revenue,
			coalesce(avg(b.price_paid) filter (where b.status='completed'),0) as avg_price
		from public.bookings b
		join public.profiles p on p.id = b.master_id
		where b.starts_at >= $1 and b.starts_at < $2
		group by p.full_name
		order by revenue desc
	`, from, to)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []MasterStat
	for rows.Next() {
		var m MasterStat
		if err := rows.Scan(&m.MasterName, &m.Total, &m.Completed, &m.Revenue, &m.AvgPrice); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *Repo) GetHourlyLoad(ctx context.Context, from, to time.Time) ([]HourlyLoad, error) {
	rows, err := r.db.Query(ctx, `
		select
			extract(hour from starts_at)::int,
			extract(dow  from starts_at)::int,
			count(*)
		from public.bookings
		where starts_at >= $1 and starts_at < $2
		  and status not in ('cancelled')
		group by 1,2 order by 2,1
	`, from, to)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []HourlyLoad
	for rows.Next() {
		var h HourlyLoad
		if err := rows.Scan(&h.HourOfDay, &h.DayOfWeek, &h.Count); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, rows.Err()
}