package config

import (
	"log"
	"os"
)

type Config struct {
	Port          string
	SupabaseURL   string
	SupabaseKey   string
	DatabaseURL   string
	JWTPublicKeyX string
	JWTPublicKeyY string
	ResendAPIKey  string
	EmailFrom     string
	GeminiAPIKey  string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		SupabaseURL:   mustEnv("SUPABASE_URL"),
		SupabaseKey:   mustEnv("SUPABASE_SERVICE_KEY"),
		DatabaseURL:   mustEnv("DATABASE_URL"),
		JWTPublicKeyX: mustEnv("JWT_PUBLIC_KEY_X"),
		JWTPublicKeyY: mustEnv("JWT_PUBLIC_KEY_Y"),
		ResendAPIKey:  getEnv("RESEND_API_KEY", ""),
		EmailFrom:     getEnv("EMAIL_FROM", "booking.servive@"),
		GeminiAPIKey:  getEnv("GEMINI_API_KEY", ""),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env %s is required", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
