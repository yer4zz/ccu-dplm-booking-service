package reminder

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"booking-service/pkg/notify"
)

type ReminderJob struct {
	db       *pgxpool.Pool
	notifier *notify.EmailNotifier
}

func New(db *pgxpool.Pool, n *notify.EmailNotifier) *ReminderJob {
	return &ReminderJob{db: db, notifier: n}
}

func (r *ReminderJob) RunAll() {
	r.SendUpcomingReminders()
	r.SendPersonalizedReminders()
	r.ExpireWaitlist()
}

func (r *ReminderJob) SendUpcomingReminders() {
	ctx := context.Background()

	rows, err := r.db.Query(ctx, `
		select
			b.id,
			cp.full_name  as client_name,
			cu.email      as client_email,
			mp.full_name  as master_name,
			s.name        as service_name,
			b.starts_at,
			s.duration_min
		from public.bookings b
		join public.profiles cp on cp.id = b.client_id
		join auth.users     cu on cu.id = b.client_id
		join public.profiles mp on mp.id = b.master_id
		join public.services  s on  s.id = b.service_id
		where b.status in ('pending', 'confirmed')
		  and b.starts_at between now() + interval '23 hours'
		                      and now() + interval '25 hours'
		  and not exists (
		    select 1 from public.reminder_log rl
		    where rl.booking_id = b.id and rl.type = 'upcoming_24h'
		  )
	`)
	if err != nil {
		fmt.Printf("[reminder] upcoming query error: %v\n", err)
		return
	}
	defer rows.Close()

	type row struct {
		ID          string
		ClientName  string
		ClientEmail string
		MasterName  string
		ServiceName string
		StartsAt    time.Time
		DurationMin int
	}

	var upcoming []row
	for rows.Next() {
		var x row
		if err := rows.Scan(
			&x.ID, &x.ClientName, &x.ClientEmail,
			&x.MasterName, &x.ServiceName, &x.StartsAt, &x.DurationMin,
		); err != nil {
			continue
		}
		upcoming = append(upcoming, x)
	}

	for _, u := range upcoming {
		err := r.notifier.SendUpcomingReminder(ctx, &notify.ReminderData{
			ClientName:  u.ClientName,
			ClientEmail: u.ClientEmail,
			MasterName:  u.MasterName,
			ServiceName: u.ServiceName,
			StartsAt:    u.StartsAt,
			DurationMin: u.DurationMin,
		})
		if err != nil {
			fmt.Printf("[reminder] ❌ ошибка отправки: %v\n", err)
			continue
		}
		r.db.Exec(ctx, `
			insert into public.reminder_log (booking_id, type)
			values ($1, 'upcoming_24h')
		`, u.ID)
		fmt.Printf("[reminder] напоминание → %s (%s)\n", u.ClientEmail, u.ServiceName)
	}
}

func (r *ReminderJob) SendPersonalizedReminders() {
	ctx := context.Background()

	rows, err := r.db.Query(ctx, `
		with visits as (
		  select
		    client_id,
		    starts_at::date as visit_date,
		    lag(starts_at::date) over (
		      partition by client_id order by starts_at
		    ) as prev_visit_date
		  from public.bookings
		  where status = 'completed'
		),
		intervals as (
		  select
		    client_id,
		    avg(visit_date - prev_visit_date)::int as avg_interval_days,
		    max(visit_date) as last_visit,
		    count(*) as visit_count
		  from visits
		  where prev_visit_date is not null
		  group by client_id
		  having count(*) >= 1
		)
		select
		  i.client_id::text,
		  p.full_name,
		  u.email,
		  i.avg_interval_days,
		  i.last_visit
		from intervals i
		join public.profiles p on p.id = i.client_id
		join auth.users     u on u.id  = i.client_id
		where i.avg_interval_days is not null
		  and i.last_visit + (i.avg_interval_days - 7) * interval '1 day' <= now()
		  and i.last_visit + i.avg_interval_days * interval '1 day'        >= now()
		  and not exists (
		    select 1 from public.reminder_log rl
		    where rl.client_id = i.client_id
		      and rl.type = 'personalized'
		      and rl.created_at > now() - interval '7 days'
		  )
	`)
	if err != nil {
		fmt.Printf("[reminder] personalized query error: %v\n", err)
		return
	}
	defer rows.Close()

	type clientRow struct {
		ClientID        string
		FullName        string
		Email           string
		AvgIntervalDays int
		LastVisit       time.Time
	}

	var clients []clientRow
	for rows.Next() {
		var c clientRow
		if err := rows.Scan(
			&c.ClientID, &c.FullName, &c.Email,
			&c.AvgIntervalDays, &c.LastVisit,
		); err != nil {
			continue
		}
		clients = append(clients, c)
	}

	for _, c := range clients {
		err := r.notifier.SendPersonalizedReminder(ctx, &notify.PersonalizedData{
			ClientName:      c.FullName,
			ClientEmail:     c.Email,
			AvgIntervalDays: c.AvgIntervalDays,
			LastVisit:       c.LastVisit,
		})
		if err != nil {
			fmt.Printf("[reminder] персон. ошибка: %v\n", err)
			continue
		}
		r.db.Exec(ctx, `
			insert into public.reminder_log (client_id, type)
			values ($1, 'personalized')
		`, c.ClientID)
		fmt.Printf("[reminder] ✅ персон. → %s (каждые %d дней)\n",
			c.Email, c.AvgIntervalDays)
	}
}

func (r *ReminderJob) ExpireWaitlist() {
	_, err := r.db.Exec(context.Background(), `
		update public.waitlist
		set status = 'expired'
		where status = 'notified'
		  and expires_at < now()
	`)
	if err != nil {
		fmt.Printf("[reminder] expire waitlist error: %v\n", err)
	}
}