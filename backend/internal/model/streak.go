package model

import (
	"time"
	"github.com/google/uuid"
)

type BeautyStreak struct {
	ClientID      uuid.UUID `db:"client_id"      json:"client_id"`
	CurrentStreak int       `db:"current_streak" json:"current_streak"`
	LongestStreak int       `db:"longest_streak" json:"longest_streak"`
	LastVisitDate *string   `db:"last_visit_date" json:"last_visit_date"`
	StreakDeadline *string  `db:"streak_deadline" json:"streak_deadline"`
	Level         string    `db:"level"          json:"level"`
	UpdatedAt     time.Time `db:"updated_at"     json:"updated_at"`

	DaysUntilDeadline int     `db:"-" json:"days_until_deadline"`
	NextLevelAt       int     `db:"-" json:"next_level_at"`
	Progress          float64 `db:"-" json:"progress"`
	StreakDiscount    float64 `db:"-" json:"streak_discount"`
}

type SOSRequest struct {
	ID                 uuid.UUID `db:"id"                   json:"id"`
	ClientID           uuid.UUID `db:"client_id"            json:"client_id"`
	MasterID           uuid.UUID `db:"master_id"            json:"master_id"`
	ServiceID          uuid.UUID `db:"service_id"           json:"service_id"`
	PreferredRangeStart time.Time `db:"preferred_range_start" json:"preferred_range_start"`
	PreferredRangeEnd   time.Time `db:"preferred_range_end"   json:"preferred_range_end"`
	BasePrice          float64   `db:"base_price"           json:"base_price"`
	SOSPrice           float64   `db:"sos_price"            json:"sos_price"`
	Status             string    `db:"status"               json:"status"`
	ClientNote         string    `db:"client_note"          json:"client_note"`
	MasterNote         string    `db:"master_note"          json:"master_note"`
	ExpiresAt          time.Time `db:"expires_at"           json:"expires_at"`
	CreatedAt          time.Time `db:"created_at"           json:"created_at"`

	ClientName  string `db:"-" json:"client_name,omitempty"`
	ServiceName string `db:"-" json:"service_name,omitempty"`
}

type CreateSOSRequest struct {
	MasterID           uuid.UUID `json:"master_id"            binding:"required"`
	ServiceID          uuid.UUID `json:"service_id"           binding:"required"`
	PreferredStart     time.Time `json:"preferred_start"      binding:"required"`
	PreferredEnd       time.Time `json:"preferred_end"        binding:"required"`
	ClientNote         string    `json:"client_note"`
}