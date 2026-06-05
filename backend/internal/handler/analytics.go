package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"booking-service/internal/repository"
	"booking-service/internal/service"
)

func GetAnalytics(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		from, to, forecastDays := parseDateRange(c)

		overview, err := repo.GetOverviewStats(c.Request.Context(), from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		daily, err := repo.GetRevenueByDay(c.Request.Context(), from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		services, err := repo.GetTopServices(c.Request.Context(), from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		masters, err := repo.GetMasterStats(c.Request.Context(), from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		revenues := make([]float64, len(daily))
		days     := make([]time.Time, len(daily))
		for i, d := range daily {
			revenues[i] = d.Revenue
			days[i]     = d.Day
		}

		forecast := service.HoltWinters(revenues, days, forecastDays)

		c.JSON(http.StatusOK, gin.H{
			"overview": overview,
			"daily":    daily,
			"forecast": forecast,
			"services": services,
			"masters":  masters,
			"from":     from.Format("2006-01-02"),
			"to":       to.Format("2006-01-02"),
		})
	}
}

func parseDateRange(c *gin.Context) (time.Time, time.Time, int) {
	rangeParam := c.DefaultQuery("range", "month")
	now        := time.Now()

	switch rangeParam {
	case "week":
		return now.AddDate(0, 0, -7), now, 7
	case "quarter":
		return now.AddDate(0, -3, 0), now, 30
	case "year":
		return now.AddDate(-1, 0, 0), now, 90
	default:
		return now.AddDate(0, -1, 0), now, 14
	}
}