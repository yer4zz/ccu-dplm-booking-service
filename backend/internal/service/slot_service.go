package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"booking-service/internal/model"
	"booking-service/internal/repository"
)

type SlotService struct {
	Repo *repository.Repo
}

func NewSlotService(repo *repository.Repo) *SlotService {
	return &SlotService{Repo: repo}
}

func (s *SlotService) GetAvailableSlots(
	ctx context.Context,
	masterID   uuid.UUID,
	serviceIDs []uuid.UUID,
	date       time.Time,
) ([]model.Slot, error) {

	dow := int(date.Weekday())

	schedule, err := s.Repo.GetDaySchedule(ctx, masterID.String(), dow)
	if err != nil || schedule == nil {
		return []model.Slot{}, nil
	}

	ids := make([]string, len(serviceIDs))
	for i, id := range serviceIDs {
		ids[i] = id.String()
	}

	var totalDuration int
	err = s.Repo.DB().QueryRow(ctx, `
		SELECT COALESCE(SUM(duration_min), 0)
		FROM public.services
		WHERE id = ANY($1::uuid[])
		  AND is_active = true
	`, ids).Scan(&totalDuration)
	if err != nil || totalDuration == 0 {
		return nil, fmt.Errorf("cannot get duration: %w", err)
	}

	duration := time.Duration(totalDuration) * time.Minute

	booked, err := s.Repo.GetBookingsForDay(ctx, masterID, date)
	if err != nil {
		return nil, err
	}

	almatyZone := time.FixedZone("UTC+5", 5*60*60)

	parseT := func(t string) time.Time {
	    var h, m int
	    fmt.Sscanf(t, "%d:%d", &h, &m)
	    return time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, almatyZone)
	}
	
	cur    := parseT(schedule.StartTime)
	dayEnd := parseT(schedule.EndTime)
	now    := time.Now().In(almatyZone)
	step   := 30 * time.Minute
	
	minTime := now.Add(30 * time.Minute)
	
	var slots []model.Slot
	for !cur.Add(duration).After(dayEnd) {
	    end := cur.Add(duration)
	    if cur.After(minTime) {
	        if !hasOverlap(cur, end, booked) {
	            slots = append(slots, model.Slot{
	                StartsAt: cur,
	                EndsAt:   end,
	            })
	        }
	    }
	    cur = cur.Add(step)
	}

	if slots == nil {
		return []model.Slot{}, nil
	}
	return slots, nil
}

func hasOverlap(start, end time.Time, bookings []model.Booking) bool {
	for _, b := range bookings {
		if start.Before(b.EndsAt) && end.After(b.StartsAt) {
			return true
		}
	}
	return false
}

