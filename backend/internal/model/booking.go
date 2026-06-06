package model

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
	StatusNoShow    BookingStatus = "no_show"
)

type Booking struct {
	ID          uuid.UUID     `db:"id"         json:"id"`
	ClientID    uuid.UUID     `db:"client_id"  json:"client_id"`
	MasterID    uuid.UUID     `db:"master_id"  json:"master_id"`
	ServiceID   uuid.UUID     `db:"service_id" json:"service_id"`
	StartsAt    time.Time     `db:"starts_at"  json:"starts_at"`
	EndsAt      time.Time     `db:"ends_at"    json:"ends_at"`
	Status      BookingStatus `db:"status"     json:"status"`
	PricePaid   float64       `db:"price_paid" json:"price_paid"`
	Notes       string        `db:"notes"      json:"notes,omitempty"`
	VibeMode    string        `db:"vibe_mode"  json:"vibe_mode,omitempty"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
	ClientEmail string        `db:"client_email" json:"client_email,omitempty"`
	ServiceName string        `db:"-"            json:"service_name,omitempty"`
	MasterName  string        `db:"-"            json:"master_name,omitempty"`
	ClientName  string        `db:"-"            json:"client_name,omitempty"`
}

type CreateBookingRequest struct {
	MasterID   uuid.UUID   `json:"master_id"   binding:"required"`
	ServiceIDs []uuid.UUID `json:"service_ids" binding:"required,min=1"`
	StartsAt   time.Time   `json:"starts_at"   binding:"required"`
	Notes      string      `json:"notes"`
	PointsUsed int         `json:"points_used"`
	
	VibeMode    string      `json:"vibe_mode"`
	VibeNote    string      `json:"vibe_note"`
}

type Slot struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type BookingNotification struct {
	BookingID   string
	ClientName  string
	ClientEmail string
	MasterName  string
	MasterEmail string
	ServiceName string
	StartsAt    time.Time
	DurationMin int
	Price       float64
}


type RescheduleRequest struct {
	NewMasterID *uuid.UUID `json:"new_master_id"`
	NewStartsAt time.Time  `json:"new_starts_at"  binding:"required"`
}