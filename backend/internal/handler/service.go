package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"booking-service/internal/cache"
	"booking-service/internal/repository"
)

func ListServices(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		services, err := repo.ListServices(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}
		c.JSON(http.StatusOK, services)
	}
}

func ListServicesCached(repo *repository.Repo, c *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if cached, ok := c.Get("services"); ok {
			ctx.JSON(http.StatusOK, cached)
			return
		}
		services, err := repo.ListServices(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}
		c.Set("services", services, 5*time.Minute)
		ctx.JSON(http.StatusOK, services)
	}
}