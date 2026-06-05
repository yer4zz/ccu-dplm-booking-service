package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name"`
}

func Register(supabaseURL, supabaseKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, _ := json.Marshal(map[string]any{
			"email":    req.Email,
			"password": req.Password,
			"data":     map[string]string{"full_name": req.FullName, "role": "client"},
		})
		resp, err := supabasePost(supabaseURL+"/auth/v1/signup", supabaseKey, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
			return
		}
		c.Data(resp.StatusCode, "application/json", mustReadBody(resp))
	}
}

func Login(supabaseURL, supabaseKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, _ := json.Marshal(map[string]string{
			"email": req.Email, "password": req.Password,
		})
		resp, err := supabasePost(
			supabaseURL+"/auth/v1/token?grant_type=password", supabaseKey, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
			return
		}
		c.Data(resp.StatusCode, "application/json", mustReadBody(resp))
	}
}

func GetMe(supabaseURL, supabaseKey string, pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		role   := c.GetString("role")

		var fullName, phone, email string
		pool.QueryRow(c.Request.Context(), `
			SELECT
				coalesce(p.full_name, ''),
				coalesce(p.phone, ''),
				coalesce(u.email, '')
			FROM public.profiles p
			JOIN auth.users u ON u.id = p.id
			WHERE p.id = $1::uuid
		`, userID).Scan(&fullName, &phone, &email)

		c.JSON(http.StatusOK, gin.H{
			"id":                 userID,
			"email":              email,
			"profile_full_name":  fullName,
			"profile_phone":      phone,
			"app_metadata": map[string]any{
				"role": role,
			},
		})
	}
}

func UpdateProfile(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := uuid.Parse(c.GetString("user_id"))

		var body struct {
			FullName string `json:"full_name"`
			Phone    string `json:"phone"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := pool.Exec(context.Background(), `
			update public.profiles
			set full_name = coalesce(nullif($1, ''), full_name),
			    phone     = coalesce(nullif($2, ''), phone)
			where id = $3
		`, body.FullName, body.Phone, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func supabasePost(url, key string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("apikey", key)
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func mustReadBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return b
}

func trimBearer(h string) string {
	return strings.TrimPrefix(h, "Bearer ")
}