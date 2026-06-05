package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"booking-service/internal/model"
	"booking-service/internal/repository"
	"booking-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBooking(svc *service.BookingService) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientID, _ := uuid.Parse(c.GetString("user_id"))

        bodyBytes, _ := io.ReadAll(c.Request.Body)
		fmt.Printf("[create-booking] body: %s\n", string(bodyBytes))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req model.CreateBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("[create-booking] bind error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

        booking, err := svc.Create(c.Request.Context(), clientID, req)
        if err != nil {
            switch err {
            case service.ErrSlotUnavailable:
                c.JSON(http.StatusConflict, gin.H{"error": "slot_unavailable"})
            case service.ErrServiceMismatch:
                c.JSON(http.StatusBadRequest, gin.H{"error": "service_not_provided"})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            }
            return
        }
        c.JSON(http.StatusCreated, booking)
    }
}

func MyBookings(repo *repository.Repo) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientID, _ := uuid.Parse(c.GetString("user_id"))
        list, err := repo.GetClientBookings(c.Request.Context(), clientID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
            return
        }
        c.JSON(http.StatusOK, list)
    }
}

func CancelBooking(svc *service.BookingService) gin.HandlerFunc {
    return func(c *gin.Context) {
        bookingID, err := uuid.Parse(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }
        requesterID, _ := uuid.Parse(c.GetString("user_id"))

        if err := svc.Cancel(c.Request.Context(), bookingID, requesterID); err != nil {
            switch err {
            case service.ErrForbidden:
                c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
            case service.ErrAlreadyDone:
                c.JSON(http.StatusBadRequest, gin.H{"error": "already_completed"})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            }
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
    }
}

func MasterBookings(repo *repository.Repo) gin.HandlerFunc {
    return func(c *gin.Context) {
        masterID, _ := uuid.Parse(c.GetString("user_id"))
        list, err := repo.GetMasterBookings(c.Request.Context(), masterID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
            return
        }
        c.JSON(http.StatusOK, list)
    }
}

func UpdateBookingStatus(repo *repository.Repo) gin.HandlerFunc {
    return func(c *gin.Context) {
        bookingID, err := uuid.Parse(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }

        var body struct {
            Status model.BookingStatus `json:"status" binding:"required"`
        }
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := repo.UpdateStatus(c.Request.Context(), bookingID, body.Status); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": body.Status})
    }
}

func parseDate(s string) (time.Time, error) {
    return time.Parse("2006-01-02", s)
}



func RescheduleBooking(svc *service.BookingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookingID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		requesterID, _ := uuid.Parse(c.GetString("user_id"))

		var req model.RescheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := svc.Reschedule(c.Request.Context(), bookingID, requesterID, req)
		if err != nil {
			switch err {
			case service.ErrSlotUnavailable:
				c.JSON(http.StatusConflict, gin.H{"error": "slot_unavailable"})
			case service.ErrForbidden:
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, updated)
	}
}


func GetAutoReschedules(repo *repository.Repo) gin.HandlerFunc {
    return func(c *gin.Context) {
        userIDStr := c.GetString("user_id")
        clientID, err := uuid.Parse(userIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
            return
        }

        list, err := repo.GetAutoReschedules(c.Request.Context(), clientID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if list == nil {
            list = []repository.AutoReschedule{}
        }

        c.JSON(http.StatusOK, list)
    }
}

func DeclineAutoReschedule(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		id       := c.Param("id")
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		if err := repo.DeclineAutoReschedule(c.Request.Context(), id, clientID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func GetAvailableSlots(svc *service.SlotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID, err := uuid.Parse(c.Query("master_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid master_id"})
			return
		}

		almatyZone := time.FixedZone("UTC+5", 5*60*60)
		date, err := time.ParseInLocation("2006-01-02", c.Query("date"), almatyZone)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}

		serviceIDsStr := c.QueryArray("service_ids")
		if len(serviceIDsStr) == 0 {
			if s := c.Query("service_id"); s != "" {
				serviceIDsStr = []string{s}
			}
		}
		if len(serviceIDsStr) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "service_id required"})
			return
		}

		var serviceIDs []uuid.UUID
		for _, s := range serviceIDsStr {
			id, err := uuid.Parse(s)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id"})
				return
			}
			serviceIDs = append(serviceIDs, id)
		}

		slots, err := svc.GetAvailableSlots(c.Request.Context(), masterID, serviceIDs, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"slots": slots})
	}
}

func AdminBookings(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := repo.DB().Query(c.Request.Context(), `
			SELECT
				b.id::text,
				b.client_id::text,
				b.master_id::text,
				b.service_id::text,
				b.starts_at,
				b.ends_at,
				b.status,
				b.price_paid,
				coalesce(b.notes, ''),
				coalesce(p_client.full_name, '') as client_name,
				coalesce(p_master.full_name, '') as master_name,
				coalesce(s.name, '')             as service_name
			FROM public.bookings b
			LEFT JOIN public.profiles p_client ON p_client.id = b.client_id
			LEFT JOIN public.profiles p_master ON p_master.id = b.master_id
			LEFT JOIN public.services s        ON s.id        = b.service_id
			ORDER BY b.starts_at DESC
			LIMIT 200
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type bookingRow struct {
			ID          string    `json:"id"`
			ClientID    string    `json:"client_id"`
			MasterID    string    `json:"master_id"`
			ServiceID   string    `json:"service_id"`
			StartsAt    time.Time `json:"starts_at"`
			EndsAt      time.Time `json:"ends_at"`
			Status      string    `json:"status"`
			PricePaid   float64   `json:"price_paid"`
			Notes       string    `json:"notes"`
			ClientName  string    `json:"client_name"`
			MasterName  string    `json:"master_name"`
			ServiceName string    `json:"service_name"`
		}

		var list []bookingRow
		for rows.Next() {
			var b bookingRow
			if err := rows.Scan(
				&b.ID, &b.ClientID, &b.MasterID, &b.ServiceID,
				&b.StartsAt, &b.EndsAt, &b.Status, &b.PricePaid, &b.Notes,
				&b.ClientName, &b.MasterName, &b.ServiceName,
			); err != nil {
				continue
			}
			list = append(list, b)
		}
		if list == nil {
			list = []bookingRow{}
		}
		c.JSON(http.StatusOK, list)
	}
}