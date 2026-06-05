package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"booking-service/internal/model"
	"booking-service/internal/repository"
	"booking-service/pkg/notify"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetStreak(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		streak, err := repo.GetBeautyStreak(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, streak)
	}
}

func CreateSOS(repo *repository.Repo, notifier *notify.EmailNotifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))

		var req model.CreateSOSRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var basePrice float64
		repo.DB().QueryRow(c.Request.Context(), `
			select coalesce(ms.custom_price, s.price)
			from public.services s
			left join public.master_services ms
			  on ms.service_id = s.id and ms.master_id = $1
			where s.id = $2
		`, req.MasterID, req.ServiceID).Scan(&basePrice)

		sosReq := model.SOSRequest{
			ClientID:            clientID,
			MasterID:            req.MasterID,
			ServiceID:           req.ServiceID,
			PreferredRangeStart: req.PreferredStart,
			PreferredRangeEnd:   req.PreferredEnd,
			BasePrice:           basePrice,
			SOSPrice:            basePrice * 1.30,
			ClientNote:          req.ClientNote,
		}

		created, err := repo.CreateSOSRequest(c.Request.Context(), sosReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		go func() {
			var masterEmail, masterName, clientName string
			repo.DB().QueryRow(context.Background(), `
				select u.email, p.full_name
				from auth.users u join public.profiles p on p.id = u.id
				where u.id = $1
			`, req.MasterID).Scan(&masterEmail, &masterName)

			repo.DB().QueryRow(context.Background(), `
				select full_name from public.profiles where id = $1
			`, clientID).Scan(&clientName)

			if masterEmail != "" {
				notifier.SendSOSToMaster(context.Background(), &notify.SOSNotification{
					MasterEmail: masterEmail,
					MasterName:  masterName,
					ClientName:  clientName,
					SOSPrice:    created.SOSPrice,
					BasePrice:   created.BasePrice,
					PrefStart:   created.PreferredRangeStart,
					PrefEnd:     created.PreferredRangeEnd,
					ClientNote:  created.ClientNote,
					SOSID:       created.ID.String(),
				})
			}
		}()

		c.JSON(http.StatusCreated, created)
	}
}

func GetMasterSOS(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID, _ := uuid.Parse(c.GetString("user_id"))
		list, err := repo.GetMasterSOSRequests(c.Request.Context(), masterID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if list == nil { list = []model.SOSRequest{} }
		c.JSON(http.StatusOK, list)
	}
}

func RespondSOS(repo *repository.Repo, notifier *notify.EmailNotifier) gin.HandlerFunc {
    return func(c *gin.Context) {
        masterID, _ := uuid.Parse(c.GetString("user_id"))
        sosID, err  := uuid.Parse(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }

        var body struct {
            Accept bool   `json:"accept"`
            Note   string `json:"note"`
            StartsAt *time.Time `json:"starts_at"`
        }
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        sos, err := repo.GetSOSRequest(c.Request.Context(), sosID)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "sos not found"})
            return
        }

        if err := repo.RespondToSOS(c.Request.Context(), sosID, masterID, body.Accept, body.Note); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var createdBooking *model.Booking

        if body.Accept && body.StartsAt != nil {
            var duration int
            repo.DB().QueryRow(c.Request.Context(),
                `select duration_min from public.services where id = $1`,
                sos.ServiceID).Scan(&duration)

            endsAt := body.StartsAt.Add(time.Duration(duration) * time.Minute)

            newBooking := model.Booking{
                ClientID:  sos.ClientID,
                MasterID:  sos.MasterID,
                ServiceID: sos.ServiceID,
                StartsAt:  *body.StartsAt,
                EndsAt:    endsAt,
                Status:    model.StatusConfirmed,
                PricePaid: sos.SOSPrice,
                Notes:     fmt.Sprintf("SOS запись. %s", sos.ClientNote),
            }

            created, err := repo.CreateBooking(c.Request.Context(), newBooking)
            if err != nil {
                fmt.Printf("[sos] ошибка создания booking: %v\n", err)
            } else {
                createdBooking = created
            }
        }

        go func() {
            var clientEmail, clientName string
            repo.DB().QueryRow(context.Background(), `
                select u.email, p.full_name
                from auth.users u join public.profiles p on p.id = u.id
                where u.id = $1
            `, sos.ClientID).Scan(&clientEmail, &clientName)

            if clientEmail != "" {
                notifier.SendSOSResponse(context.Background(), &notify.SOSResponse{
                    ClientEmail: clientEmail,
                    ClientName:  clientName,
                    Accepted:    body.Accept,
                    MasterNote:  body.Note,
                    SOSPrice:    sos.SOSPrice,
                    PrefStart:   sos.PreferredRangeStart,
                })
            }
        }()

        resp := gin.H{"ok": true, "accepted": body.Accept}
        if createdBooking != nil {
            resp["booking"] = createdBooking
        }
        c.JSON(http.StatusOK, resp)
    }
}

func GetClientSOS(repo *repository.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := uuid.Parse(c.GetString("user_id"))
		list, err := repo.GetClientSOSRequests(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if list == nil { list = []model.SOSRequest{} }
		c.JSON(http.StatusOK, list)
	}
}