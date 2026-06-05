package model

import (
	"time"
	"github.com/google/uuid"
)

type LoyaltyAccount struct {
	ClientID     uuid.UUID `db:"client_id"    json:"client_id"`
	Balance      int       `db:"balance"       json:"balance"`
	TotalEarned  int       `db:"total_earned"  json:"total_earned"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

type PointTransaction struct {
	ID          uuid.UUID  `db:"id"          json:"id"`
	ClientID    uuid.UUID  `db:"client_id"   json:"client_id"`
	BookingID   *uuid.UUID `db:"booking_id"  json:"booking_id,omitempty"`
	Type        string     `db:"type"        json:"type"`
	Amount      int        `db:"amount"      json:"amount"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at"  json:"created_at"`
}

type Discount struct {
	ID         uuid.UUID `db:"id"          json:"id"`
	BookingID  uuid.UUID `db:"booking_id"  json:"booking_id"`
	Type       string    `db:"type"        json:"type"`
	Percent    float64   `db:"percent"     json:"percent"`
	PointsUsed int       `db:"points_used" json:"points_used"`
	Amount     float64   `db:"amount"      json:"amount"`
}

type WaitlistEntry struct {
	ID            uuid.UUID  `db:"id"             json:"id"`
	ClientID      uuid.UUID  `db:"client_id"      json:"client_id"`
	MasterID      uuid.UUID  `db:"master_id"      json:"master_id"`
	ServiceID     uuid.UUID  `db:"service_id"     json:"service_id"`
	PreferredDate *time.Time `db:"preferred_date" json:"preferred_date,omitempty"`
	Status        string     `db:"status"         json:"status"`
	CreatedAt     time.Time  `db:"created_at"     json:"created_at"`
}

type RescheduleOffer struct {
	ID                uuid.UUID `db:"id"                  json:"id"`
	OriginalBookingID uuid.UUID `db:"original_booking_id" json:"original_booking_id"`
	ClientID          uuid.UUID `db:"client_id"           json:"client_id"`
	Status            string    `db:"status"              json:"status"`
	OfferedAt         time.Time `db:"offered_at"          json:"offered_at"`
	ExpiresAt         time.Time `db:"expires_at"          json:"expires_at"`
}

type BookingPriceCalc struct {
	BasePrice       float64   `json:"base_price"`
	ServiceCount    int       `json:"service_count"`
	IsFirstBooking  bool      `json:"is_first_booking"`
	FirstDiscount   float64   `json:"first_discount"`
	MultiDiscount   float64   `json:"multi_discount"`
	PointsAvailable int       `json:"points_available"`
	PointsUsed      int       `json:"points_used"`
	PointsDiscount  float64   `json:"points_discount"`
	TotalDiscount   float64   `json:"total_discount"`
	FinalPrice      float64   `json:"final_price"`
}