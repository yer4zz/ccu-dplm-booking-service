package handler

import (
	"booking-service/internal/cache"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminCreateService(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Name        string  `json:"name"         binding:"required"`
			Category    string  `json:"category"     binding:"required"`
			Description string  `json:"description"`
			DurationMin int     `json:"duration_min" binding:"required"`
			Price       float64 `json:"price"        binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var id string
		err := pool.QueryRow(c.Request.Context(), `
			INSERT INTO public.services
			  (name, category, description, duration_min, price, is_active)
			VALUES ($1,$2,$3,$4,$5,true)
			RETURNING id::text
		`, body.Name, body.Category, body.Description,
			body.DurationMin, body.Price,
		).Scan(&id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cache.Global.Delete("masters")
		cache.Global.Delete("services")
		cache.Global.Delete("gallery")
		c.JSON(http.StatusCreated, gin.H{"id": id, "ok": true})
	}
}

func AdminUpdateService(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var body struct {
			Name        *string  `json:"name"`
			Category    *string  `json:"category"`
			Description *string  `json:"description"`
			DurationMin *int     `json:"duration_min"`
			Price       *float64 `json:"price"`
			IsActive    *bool    `json:"is_active"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := pool.Exec(c.Request.Context(), `
			UPDATE public.services SET
				name         = COALESCE($1, name),
				category     = COALESCE($2, category),
				description  = COALESCE($3, description),
				duration_min = COALESCE($4, duration_min),
				price        = COALESCE($5, price),
				is_active    = COALESCE($6, is_active)
			WHERE id = $7::uuid
		`, body.Name, body.Category, body.Description,
			body.DurationMin, body.Price, body.IsActive, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cache.Global.Delete("masters")
		cache.Global.Delete("services")
		cache.Global.Delete("gallery")
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func AdminDeleteService(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		pool.Exec(c.Request.Context(),
			`DELETE FROM public.master_services WHERE service_id = $1::uuid`, id)

		_, err := pool.Exec(c.Request.Context(),
			`DELETE FROM public.services WHERE id = $1::uuid`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cache.Global.Delete("masters")
		cache.Global.Delete("services")
		cache.Global.Delete("gallery")
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func AdminCreateMaster(pool *pgxpool.Pool, supabaseURL, supabaseKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email           string `json:"email"            binding:"required"`
			Password        string `json:"password"         binding:"required"`
			FullName        string `json:"full_name"        binding:"required"`
			Bio             string `json:"bio"`
			ExperienceYears int    `json:"experience_years"`
			Instagram       string `json:"instagram"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payload, _ := json.Marshal(map[string]any{
			"email":         body.Email,
			"password":      body.Password,
			"email_confirm": true,
			"user_metadata": map[string]any{
				"full_name": body.FullName,
				"role":      "master",
			},
			"app_metadata": map[string]any{
				"role": "master",
			},
		})

		req, _ := http.NewRequestWithContext(
			c.Request.Context(),
			"POST",
			supabaseURL+"/auth/v1/admin/users",
			bytes.NewReader(payload),
		)
		req.Header.Set("apikey", supabaseKey)
		req.Header.Set("Authorization", "Bearer "+supabaseKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "supabase request failed"})
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			c.JSON(http.StatusBadRequest, gin.H{"error": string(respBody)})
			return
		}

		var authResp struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(respBody, &authResp); err != nil || authResp.ID == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot parse user id"})
			return
		}

		userID := authResp.ID

		updatePayload, _ := json.Marshal(map[string]any{
		    "app_metadata": map[string]any{
		        "role": "master",
		    },
		})
		updateReq, _ := http.NewRequestWithContext(
		    c.Request.Context(),
		    "PUT",
		    supabaseURL+"/auth/v1/admin/users/"+userID,
		    bytes.NewReader(updatePayload),
		)
		updateReq.Header.Set("apikey", supabaseKey)
		updateReq.Header.Set("Authorization", "Bearer "+supabaseKey)
		updateReq.Header.Set("Content-Type", "application/json")
		http.DefaultClient.Do(updateReq)

		_, err = pool.Exec(c.Request.Context(), `
			INSERT INTO public.profiles (id, full_name, role)
			VALUES ($1::uuid, $2, 'master')
			ON CONFLICT (id) DO UPDATE SET full_name = $2, role = 'master'
		`, userID, body.FullName)
		if err != nil {
			fmt.Printf("[admin] profile insert error: %v\n", err)
		}

		_, err = pool.Exec(c.Request.Context(), `
			INSERT INTO public.masters (id, bio, experience_years, instagram, is_active)
			VALUES ($1::uuid, $2, $3, $4, true)
			ON CONFLICT (id) DO UPDATE
			  SET bio = $2, experience_years = $3, instagram = $4
		`, userID, body.Bio, body.ExperienceYears, body.Instagram)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "masters insert: " + err.Error()})
			return
		}
		cache.Global.Delete("masters")
		cache.Global.Delete("services")
		cache.Global.Delete("gallery")

		c.JSON(http.StatusCreated, gin.H{
			"ok":      true,
			"user_id": userID,
			"email":   body.Email,
		})
	}
}

func AdminUpdateMaster(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var body struct {
			Bio             *string `json:"bio"`
			ExperienceYears *int    `json:"experience_years"`
			Instagram       *string `json:"instagram"`
			IsActive        *bool   `json:"is_active"`
			FullName        *string `json:"full_name"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := pool.Exec(c.Request.Context(), `
			UPDATE public.masters SET
				bio              = COALESCE($1, bio),
				experience_years = COALESCE($2, experience_years),
				instagram        = COALESCE($3, instagram),
				is_active        = COALESCE($4, is_active)
			WHERE id = $5::uuid
		`, body.Bio, body.ExperienceYears, body.Instagram, body.IsActive, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if body.FullName != nil {
			pool.Exec(c.Request.Context(),
				`UPDATE public.profiles SET full_name = $1 WHERE id = $2::uuid`,
				*body.FullName, id)
		}
		cache.Global.Delete("masters")
		cache.Global.Delete("services")
		cache.Global.Delete("gallery")
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func AdminDeleteMaster(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := pool.Exec(c.Request.Context(),
			`UPDATE public.masters SET is_active = false WHERE id = $1::uuid`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cache.Global.Delete("masters")
		cache.Global.Delete("services")
		cache.Global.Delete("gallery")
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}


func MasterAddService(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID := c.GetString("user_id")
		var body struct {
			ServiceID   string   `json:"service_id"   binding:"required"`
			CustomPrice *float64 `json:"custom_price"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := pool.Exec(c.Request.Context(), `
			INSERT INTO public.master_services (master_id, service_id, custom_price)
			VALUES ($1::uuid, $2::uuid, $3)
			ON CONFLICT (master_id, service_id) DO UPDATE
			  SET custom_price = EXCLUDED.custom_price
		`, masterID, body.ServiceID, body.CustomPrice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"ok": true})
	}
}

func MasterRemoveService(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID  := c.GetString("user_id")
		serviceID := c.Param("service_id")

		_, err := pool.Exec(c.Request.Context(), `
			DELETE FROM public.master_services
			WHERE master_id = $1::uuid AND service_id = $2::uuid
		`, masterID, serviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func MasterUpdateServicePrice(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID  := c.GetString("user_id")
		serviceID := c.Param("service_id")
		var body struct {
			CustomPrice *float64 `json:"custom_price"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := pool.Exec(c.Request.Context(), `
			UPDATE public.master_services
			SET custom_price = $1
			WHERE master_id = $2::uuid AND service_id = $3::uuid
		`, body.CustomPrice, masterID, serviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func GetAllServices(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		showAll := c.Query("show_all") == "true"

		query := `
			SELECT id::text, name, category, coalesce(description,''),
			       duration_min, price, is_active
			FROM public.services
		`
		if !showAll {
			query += " WHERE is_active = true"
		}
		query += " ORDER BY sort_order, name"

		rows, err := pool.Query(c.Request.Context(), query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type svcRow struct {
			ID          string  `json:"id"`
			Name        string  `json:"name"`
			Category    string  `json:"category"`
			Description string  `json:"description"`
			DurationMin int     `json:"duration_min"`
			Price       float64 `json:"price"`
			IsActive    bool    `json:"is_active"`
		}
		var list []svcRow
		for rows.Next() {
			var s svcRow
			rows.Scan(&s.ID, &s.Name, &s.Category, &s.Description,
				&s.DurationMin, &s.Price, &s.IsActive)
			list = append(list, s)
		}
		if list == nil {
			list = []svcRow{}
		}
		c.JSON(http.StatusOK, list)
	}
}

func GetMasterServices(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		masterID := c.GetString("user_id")
		rows, err := pool.Query(c.Request.Context(), `
			SELECT
				s.id::text,
				s.name,
				s.category,
				coalesce(s.description, ''),
				s.duration_min,
				s.price,
				s.is_active,
				(ms.master_id IS NOT NULL) as is_assigned,
				coalesce(ms.custom_price, s.price) as effective_price
			FROM public.services s
			LEFT JOIN public.master_services ms
				ON ms.service_id = s.id
				AND ms.master_id = $1::uuid
			WHERE s.is_active = true
			ORDER BY s.sort_order, s.name
		`, masterID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type svcRow struct {
			ID             string  `json:"id"`
			Name           string  `json:"name"`
			Category       string  `json:"category"`
			Description    string  `json:"description"`
			DurationMin    int     `json:"duration_min"`
			Price          float64 `json:"price"`
			IsActive       bool    `json:"is_active"`
			IsAssigned     bool    `json:"is_assigned"`
			EffectivePrice float64 `json:"effective_price"`
		}

		var list []svcRow
		for rows.Next() {
			var s svcRow
			if err := rows.Scan(
				&s.ID, &s.Name, &s.Category, &s.Description,
				&s.DurationMin, &s.Price, &s.IsActive,
				&s.IsAssigned, &s.EffectivePrice,
			); err != nil {
				continue
			}
			list = append(list, s)
		}
		if list == nil {
			list = []svcRow{}
		}
		c.JSON(http.StatusOK, list)
	}
}