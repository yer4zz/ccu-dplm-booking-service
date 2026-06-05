package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "booking-service/internal/repository"
    "time"
    "booking-service/internal/cache"
)

func ListMasters(repo *repository.Repo) gin.HandlerFunc {
    return func(c *gin.Context) {
        masters, err := repo.ListMasters(c.Request.Context())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
            return
        }
        c.JSON(http.StatusOK, masters)
    }
}

func ListMastersCached(repo *repository.Repo, c *cache.Cache) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if cached, ok := c.Get("masters"); ok {
            ctx.JSON(http.StatusOK, cached)
            return
        }
        masters, err := repo.ListMasters(ctx.Request.Context())
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
            return
        }
        c.Set("masters", masters, 5*time.Minute)
        ctx.JSON(http.StatusOK, masters)
    }
}