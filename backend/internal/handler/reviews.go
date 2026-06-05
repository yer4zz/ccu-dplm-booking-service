package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"booking-service/internal/repository"
)

func GetPublicReviews(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := pool.Query(c.Request.Context(), `
			select
				r.id::text,
				r.rating,
				coalesce(r.comment, ''),
				coalesce(p.full_name, 'Клиент') as client_name,
				coalesce(s.name, '') as service_name
			from public.reviews r
			join public.profiles p on p.id  = r.client_id
			join public.bookings b on b.id  = r.booking_id
			join public.services s on s.id  = b.service_id
			where r.is_visible = true
			  and r.comment != ''
			order by r.created_at desc
			limit 6
		`)
		if err != nil {
			c.JSON(http.StatusOK, []any{})
			return
		}
		defer rows.Close()

		type reviewRow struct {
			ID          string `json:"id"`
			Rating      int    `json:"rating"`
			Comment     string `json:"comment"`
			ClientName  string `json:"client_name"`
			ServiceName string `json:"service_name"`
		}

		var list []reviewRow
		for rows.Next() {
			var rv reviewRow
			if err := rows.Scan(
				&rv.ID, &rv.Rating, &rv.Comment,
				&rv.ClientName, &rv.ServiceName,
			); err != nil {
				continue
			}
			list = append(list, rv)
		}
		if list == nil {
			list = []reviewRow{}
		}
		c.JSON(http.StatusOK, list)
	}
}

func DeleteReview(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID, _ := uuid.Parse(c.GetString("user_id"))
		reviewID    := c.Param("id")

		if err := repo.DeleteReviewByMaster(
			c.Request.Context(), reviewID, masterID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func DeleteReviewAdmin(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewID := c.Param("id")
		if err := repo.DeleteReviewByAdmin(
			c.Request.Context(), reviewID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}