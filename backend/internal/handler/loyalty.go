package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"booking-service/internal/model"
	"booking-service/internal/repository"
	"booking-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetLoyalty(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		account, err := repo.GetLoyaltyAccount(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, account)
	}
}

func GetTransactions(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		list, err := repo.GetPointTransactions(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func CalcBookingPrice(svc *service.BookingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		fmt.Printf("[calc-price] body: %s\n", string(bodyBytes))

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req struct {
			MasterID   uuid.UUID   `json:"master_id"   binding:"required"`
			ServiceIDs []uuid.UUID `json:"service_ids" binding:"required"`
			PointsUsed int         `json:"points_used"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("[calc-price] bind error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		calc, err := svc.CalcPrice(
			c.Request.Context(), clientID, req.MasterID, req.ServiceIDs, req.PointsUsed)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, calc)
	}
}

func AddToWaitlist(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		var req struct {
			MasterID  uuid.UUID `json:"master_id"  binding:"required"`
			ServiceID uuid.UUID `json:"service_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		entry, err := repo.AddToWaitlist(c.Request.Context(), model.WaitlistEntry{
			ClientID:  clientID,
			MasterID:  req.MasterID,
			ServiceID: req.ServiceID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, entry)
	}
}

func GetWaitlist(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		list, err := repo.GetWaitlist(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func RemoveFromWaitlist(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.Parse(c.Param("id"))
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		if err := repo.RemoveFromWaitlist(c.Request.Context(), id, clientID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func GetRescheduleOffers(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		list, err := repo.GetRescheduleOffers(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}



