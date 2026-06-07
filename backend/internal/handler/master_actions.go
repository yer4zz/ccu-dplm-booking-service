package handler

import (
	"context"
	"math"
	"net/http"
	"time"

	"booking-service/internal/model"
	"booking-service/internal/repository"
	"booking-service/internal/service"
	"booking-service/pkg/notify"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func MasterRescheduleBooking(
	repo     *repository.Repo,
	svc      *service.BookingService,
	notifier *notify.EmailNotifier,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID, _ := uuid.Parse(c.GetString("user_id"))
		bookingID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var body struct {
			NewMasterID *string    `json:"new_master_id"`
			NewStartsAt time.Time  `json:"new_starts_at" binding:"required"`
			Reason      string     `json:"reason"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		original, err := repo.GetBooking(c.Request.Context(), bookingID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		if original.MasterID != masterID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		targetMasterID := masterID
		if body.NewMasterID != nil && *body.NewMasterID != "" {
			parsed, err := uuid.Parse(*body.NewMasterID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid new_master_id"})
				return
			}
			targetMasterID = parsed
		}

		var duration int
		repo.DB().QueryRow(c.Request.Context(),
			`select duration_min from public.services where id = $1`,
			original.ServiceID).Scan(&duration)

		newEndsAt := body.NewStartsAt.Add(time.Duration(duration) * time.Minute)

		available, err := repo.IsSlotAvailableExclude(
			c.Request.Context(), targetMasterID,
			body.NewStartsAt, newEndsAt, bookingID,
		)
		if err != nil || !available {
			c.JSON(http.StatusConflict, gin.H{"error": "slot_unavailable"})
			return
		}

		repo.UpdateStatus(c.Request.Context(), bookingID, model.StatusCancelled)

		newPrice := original.PricePaid * 0.70
		newBooking, err := repo.CreateBooking(c.Request.Context(), model.Booking{
			ClientID:  original.ClientID,
			MasterID:  targetMasterID,
			ServiceID: original.ServiceID,
			StartsAt:  body.NewStartsAt,
			EndsAt:    newEndsAt,
			Status:    model.StatusConfirmed,
			PricePaid: newPrice,
			Notes:     original.Notes,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		repo.CreateMasterReschedule(
			c.Request.Context(),
			bookingID, newBooking.ID,
			original.ClientID, masterID,
			body.Reason,
		)

		go func() {
			notif, err := repo.GetBookingNotification(context.Background(), newBooking.ID)
			if err == nil && notif != nil {
				notifier.BookingAutoRescheduled(context.Background(), notif)
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"ok":          true,
			"new_booking": newBooking,
		})
	}
}

func CreateReview(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		bookingID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var body struct {
			Rating  int    `json:"rating"  binding:"required,min=1,max=5"`
			Comment string `json:"comment"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		booking, err := repo.GetBooking(c.Request.Context(), bookingID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if booking.ClientID != clientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if booking.Status != model.StatusCompleted {
			c.JSON(http.StatusBadRequest, gin.H{"error": "booking_not_completed"})
			return
		}

		if err := repo.CreateReview(c.Request.Context(), model.Review{
			BookingID: bookingID,
			ClientID:  clientID,
			MasterID:  booking.MasterID,
			Rating:    body.Rating,
			Comment:   body.Comment,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"ok": true})
	}
}

func GetMasterReviews(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID, _ := uuid.Parse(c.GetString("user_id"))
		list, err := repo.GetMasterReviews(c.Request.Context(), masterID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if list == nil { list = []model.Review{} }
		c.JSON(http.StatusOK, list)
	}
}

func GetMasterStats(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID, _ := uuid.Parse(c.GetString("user_id"))
		ctx := c.Request.Context()

		now        := time.Now()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

		var (
			totalBookings     int
			completedBookings int
			totalRevenue      float64
			avgPrice          float64
			thisMonthRevenue  float64
			thisMonthBookings int
			avgRating         float64
			totalReviews      int
			topService        string
		)

		repo.DB().QueryRow(ctx, `
			select
				count(*)                                                    as total_bookings,
				count(*) filter (where status = 'completed')               as completed_bookings,
				coalesce(sum(price_paid) filter (where status = 'completed'), 0) as total_revenue,
				coalesce(avg(price_paid) filter (where status = 'completed'), 0) as avg_price,
				coalesce(sum(price_paid) filter (where status = 'completed' and starts_at >= $2), 0) as this_month_revenue,
				count(*) filter (where starts_at >= $2)                    as this_month_bookings
			from public.bookings
			where master_id = $1
		`, masterID, monthStart).Scan(
			&totalBookings, &completedBookings,
			&totalRevenue, &avgPrice,
			&thisMonthRevenue, &thisMonthBookings,
		)

		repo.DB().QueryRow(ctx, `
			select
				coalesce(avg(rating), 0) as avg_rating,
				count(*)                 as total_reviews
			from public.reviews
			where master_id = $1
		`, masterID).Scan(&avgRating, &totalReviews)

		repo.DB().QueryRow(ctx, `
			select coalesce(s.name, '')
			from public.bookings b
			join public.services s on s.id = b.service_id
			where b.master_id = $1 and b.status = 'completed'
			group by s.name
			order by count(*) desc
			limit 1
		`, masterID).Scan(&topService)

		c.JSON(http.StatusOK, gin.H{
			"total_bookings":      totalBookings,
			"completed_bookings":  completedBookings,
			"total_revenue":       totalRevenue,
			"avg_price":           avgPrice,
			"this_month_revenue":  thisMonthRevenue,
			"this_month_bookings": thisMonthBookings,
			"avg_rating":          math.Round(avgRating*10) / 10,
			"total_reviews":       totalReviews,
			"top_service":         topService,
		})
	}
}

func GetClientReschedules(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		list, err := repo.GetClientReschedules(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if list == nil { list = []model.MasterReschedule{} }
		c.JSON(http.StatusOK, list)
	}
}

func DeclineClientReschedule(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		id := c.Param("id")
		if err := repo.DeclineMasterReschedule(
			c.Request.Context(), id, clientID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}