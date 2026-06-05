package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"booking-service/internal/config"
	"booking-service/internal/handler"
	"booking-service/internal/middleware"
	"booking-service/internal/repository"
	"booking-service/internal/service"
	"booking-service/pkg/db"
	"booking-service/pkg/notify"
	"booking-service/pkg/reminder"
	"booking-service/internal/cache"

)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	if err := middleware.InitJWKS(cfg.JWTPublicKeyX, cfg.JWTPublicKeyY); err != nil {
		log.Fatalf("cannot init JWKS: %v", err)
	}

	pool       := db.NewPool(cfg.DatabaseURL)
	db.StartKeepalive(pool)
	repo       := repository.New(pool)
	notifier   := notify.NewEmailNotifier(cfg.ResendAPIKey, cfg.EmailFrom)
	bookingSvc := service.NewBookingService(repo, notifier)
	slotSvc    := service.NewSlotService(repo)

	reminderJob := reminder.New(pool, notifier)

	c := cron.New()
	c.AddFunc("0 * * * *", func() {
		log.Println("[cron] running reminders...")
		reminderJob.RunAll()
	})
	go reminderJob.RunAll()
	c.Start()
	defer c.Stop()

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"}) 
	r.Use(middleware.CORS())

	appCache := cache.New()

	api := r.Group("/api/v1")

	api.POST("/auth/register", handler.Register(cfg.SupabaseURL, cfg.SupabaseKey))
	api.POST("/auth/login",    handler.Login(cfg.SupabaseURL, cfg.SupabaseKey))
	api.GET("/services", handler.ListServicesCached(repo, appCache))
	api.GET("/masters",  handler.ListMastersCached(repo, appCache))
	api.GET("/slots",          handler.GetAvailableSlots(slotSvc))

	api.POST("/chat", handler.ChatGemini(pool, cfg.GeminiAPIKey))

	api.GET("/reviews/public", handler.GetPublicReviews(pool))

	auth := api.Group("/")
	auth.Use(middleware.JWTAuth())

	auth.GET("/me", handler.GetMe(cfg.SupabaseURL, cfg.SupabaseKey, pool))
	auth.PATCH("/profile",              handler.UpdateProfile(pool))
	auth.POST("/bookings",              handler.CreateBooking(bookingSvc))
	auth.POST("/bookings/calc-price",   handler.CalcBookingPrice(bookingSvc))
	auth.GET("/bookings/my",            handler.MyBookings(repo))
	auth.PATCH("/bookings/:id/cancel",  handler.CancelBooking(bookingSvc))

	auth.GET("/loyalty",              handler.GetLoyalty(repo))
	auth.GET("/loyalty/transactions", handler.GetTransactions(repo))

	auth.POST("/waitlist",   handler.AddToWaitlist(repo))
	auth.GET("/waitlist",    handler.GetWaitlist(repo))
	auth.DELETE("/waitlist/:id", handler.RemoveFromWaitlist(repo))

	auth.PATCH("/bookings/:id/reschedule", handler.RescheduleBooking(bookingSvc))

	auth.GET("/reschedules",          handler.GetClientReschedules(repo))
	auth.DELETE("/reschedules/:id",   handler.DeclineClientReschedule(repo))
	auth.POST("/bookings/:id/review", handler.CreateReview(repo))

	master := auth.Group("/master")
	master.Use(middleware.RequireRole("master"))
	master.GET("/bookings",       handler.MasterBookings(repo))
	master.PATCH("/bookings/:id", handler.UpdateBookingStatus(repo))

	admin := auth.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	admin.GET("/analytics",       handler.GetAnalytics(repo))

	auth.GET("/streak", handler.GetStreak(repo))

	auth.POST("/sos",    handler.CreateSOS(repo, notifier))
	auth.GET("/sos",     handler.GetClientSOS(repo))

	auth.GET("/auto-reschedules",          handler.GetAutoReschedules(repo))
	auth.DELETE("/auto-reschedules/:id",   handler.DeclineAutoReschedule(repo))

	master.GET("/sos",          handler.GetMasterSOS(repo))
	master.PATCH("/sos/:id",    handler.RespondSOS(repo, notifier))

	master.POST("/bookings/:id/reschedule", handler.MasterRescheduleBooking(repo, bookingSvc, notifier))
	master.GET("/reviews",                  handler.GetMasterReviews(repo))
	master.GET("/stats",                    handler.GetMasterStats(repo))

	api.GET("/gallery",  handler.GetGallery(pool, appCache))

	master.DELETE("/reviews/:id", handler.DeleteReview(repo))
	admin.DELETE("/reviews/:id",  handler.DeleteReviewAdmin(repo))

	master.POST("/gallery",       handler.UploadGalleryWork(pool, cfg.SupabaseURL, cfg.SupabaseKey))
	master.DELETE("/gallery/:id", handler.DeleteGalleryWork(pool))
	admin.POST("/gallery",        handler.UploadGalleryWork(pool, cfg.SupabaseURL, cfg.SupabaseKey))
	admin.DELETE("/gallery/:id",  handler.DeleteGalleryWork(pool))


	master.GET("/my-services",                   handler.GetMasterServices(pool))
	master.POST("/my-services",                  handler.MasterAddService(pool))
	master.DELETE("/my-services/:service_id",    handler.MasterRemoveService(pool))
	master.PATCH("/my-services/:service_id",     handler.MasterUpdateServicePrice(pool))

	admin.GET("/services/all",    handler.GetAllServices(pool))
	admin.POST("/services",       handler.AdminCreateService(pool))
	admin.PATCH("/services/:id",  handler.AdminUpdateService(pool))
	admin.DELETE("/services/:id", handler.AdminDeleteService(pool))
	
	admin.POST("/masters", handler.AdminCreateMaster(pool, cfg.SupabaseURL, cfg.SupabaseKey))
	admin.PATCH("/masters/:id",  handler.AdminUpdateMaster(pool))
	admin.DELETE("/masters/:id", handler.AdminDeleteMaster(pool))

	admin.GET("/bookings",       handler.AdminBookings(repo))
	admin.PATCH("/bookings/:id", handler.UpdateBookingStatus(repo))

	admin.GET("/masters",        handler.ListMasters(repo)) 

	log.Printf("server on :%s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}