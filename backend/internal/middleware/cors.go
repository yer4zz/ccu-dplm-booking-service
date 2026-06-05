package middleware

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:5173", 
            "https://yourdomain.com",
            "https://*.ngrok-free.app", 
            "https://*.ngrok-free.dev",
            "https://ccu-dplm-booking-service.vercel.app",
            "https://*.vercel.app",
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type", "ngrok-skip-browser-warning"},
        AllowCredentials: true,
        AllowWildcard:    true,
    })
}