package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"booking-service/internal/model"
)

func (r *Repo) GetBeautyStreak(ctx context.Context, clientID uuid.UUID) (*model.BeautyStreak, error) {
	var s model.BeautyStreak
	err := r.db.QueryRow(ctx, `
		select client_id, current_streak, longest_streak,
		       to_char(last_visit_date,'YYYY-MM-DD'),
		       to_char(streak_deadline,'YYYY-MM-DD'),
		       level, updated_at
		from public.beauty_streaks
		where client_id = $1
	`, clientID).Scan(
		&s.ClientID, &s.CurrentStreak, &s.LongestStreak,
		&s.LastVisitDate, &s.StreakDeadline,
		&s.Level, &s.UpdatedAt,
	)
	if err != nil {
		r.db.Exec(ctx,
			`insert into public.beauty_streaks (client_id) values ($1) on conflict do nothing`,
			clientID)
		return &model.BeautyStreak{ClientID: clientID, Level: "bronze"}, nil
	}

	if s.StreakDeadline != nil {
		deadline, _ := time.Parse("2006-01-02", *s.StreakDeadline)
		days := int(time.Until(deadline).Hours() / 24)
		if days < 0 {
			days = 0
		}
		s.DaysUntilDeadline = days
	}

	nextLevels := map[string]int{
		"bronze":   4,
		"silver":   9,
		"gold":     16,
		"platinum": 999,
	}
	s.NextLevelAt = nextLevels[s.Level]
	prevLevels := map[string]int{
		"bronze": 0, "silver": 4, "gold": 9, "platinum": 16,
	}
	prev := prevLevels[s.Level]
	next := nextLevels[s.Level]
	if next > 900 {
		s.Progress = 100
	} else {
		s.Progress = float64(s.CurrentStreak-prev) / float64(next-prev) * 100
	}

	discounts := map[string]float64{
		"bronze": 0, "silver": 3, "gold": 5, "platinum": 10,
	}
	s.StreakDiscount = discounts[s.Level]

	return &s, nil
}