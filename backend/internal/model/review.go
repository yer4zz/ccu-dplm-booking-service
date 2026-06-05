package model

import (
	"time"
	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	BookingID uuid.UUID `db:"booking_id" json:"booking_id"`
	ClientID  uuid.UUID `db:"client_id"  json:"client_id"`
	MasterID  uuid.UUID `db:"master_id"  json:"master_id"`
	Rating    int       `db:"rating"     json:"rating"`
	Comment   string    `db:"comment"    json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type MasterReschedule struct {
	ID                 uuid.UUID  `db:"id"                   json:"id"`
	OriginalBookingID  uuid.UUID  `db:"original_booking_id"  json:"original_booking_id"`
	NewBookingID       *uuid.UUID `db:"new_booking_id"       json:"new_booking_id"`
	ClientID           uuid.UUID  `db:"client_id"            json:"client_id"`
	MasterID           uuid.UUID  `db:"master_id"            json:"master_id"`
	Reason             string     `db:"reason"               json:"reason"`
	Status             string     `db:"status"               json:"status"`
	CreatedAt          time.Time  `db:"created_at"           json:"created_at"`

	ServiceName    string    `db:"-" json:"service_name,omitempty"`
	NewMasterName  string    `db:"-" json:"new_master_name,omitempty"`
	NewStartsAt    time.Time `db:"-" json:"new_starts_at,omitempty"`
	NewPrice       float64   `db:"-" json:"new_price,omitempty"`
}

type MasterStats struct {
	TotalBookings     int64   `json:"total_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	AvgRating         float64 `json:"avg_rating"`
	TotalReviews      int64   `json:"total_reviews"`
	AvgPrice          float64 `json:"avg_price"`
	TopService        string  `json:"top_service"`
	ThisMonthRevenue  float64 `json:"this_month_revenue"`
	ThisMonthBookings int64   `json:"this_month_bookings"`
}