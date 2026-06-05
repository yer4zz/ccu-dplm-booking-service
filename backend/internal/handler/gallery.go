package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"booking-service/internal/cache"
)

type GalleryWork struct {
	ID          string  `json:"id"`
	MasterID    *string `json:"master_id"`
	ServiceID   *string `json:"service_id"`
	ImageURL    string  `json:"image_url"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	MasterName  string  `json:"master_name"`
	ServiceName string  `json:"service_name"`
	CreatedAt   string  `json:"created_at"`
}

func GetGallery(pool *pgxpool.Pool, c *cache.Cache) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        masterID  := ctx.Query("master_id")
        serviceID := ctx.Query("service_id")

        if masterID == "" && serviceID == "" {
            if cached, ok := c.Get("gallery"); ok {
                ctx.JSON(http.StatusOK, cached)
                return
            }
        }

        query := `
            select
                g.id::text, g.master_id::text, g.service_id::text,
                g.image_url, coalesce(g.title,''), coalesce(g.description,''),
                coalesce(p.full_name,''), coalesce(s.name,''), g.created_at::text
            from public.gallery_works g
            left join public.profiles p on p.id = g.master_id
            left join public.services s on s.id = g.service_id
            where g.is_visible = true
        `
        args := []any{}
        if masterID != "" {
            args = append(args, masterID)
            query += fmt.Sprintf(" and g.master_id = $%d::uuid", len(args))
        }
        if serviceID != "" {
            args = append(args, serviceID)
            query += fmt.Sprintf(" and g.service_id = $%d::uuid", len(args))
        }
        query += " order by g.sort_order, g.created_at desc limit 20"

        rows, err := pool.Query(ctx.Request.Context(), query, args...)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()

        var list []GalleryWork
        for rows.Next() {
            var w GalleryWork
            if err := rows.Scan(
                &w.ID, &w.MasterID, &w.ServiceID,
                &w.ImageURL, &w.Title, &w.Description,
                &w.MasterName, &w.ServiceName, &w.CreatedAt,
            ); err != nil {
                continue
            }
            list = append(list, w)
        }
        if list == nil { list = []GalleryWork{} }

        if masterID == "" && serviceID == "" {
            c.Set("gallery", list, 5*time.Minute)
        }
        ctx.JSON(http.StatusOK, list)
    }
}

func sanitizeFileName(name string) string {
    var result strings.Builder
    for _, r := range name {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
            (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
            result.WriteRune(r)
        } else {
            result.WriteRune('_')
        }
    }
    s := result.String()
    if s == "" || s == "." {
        return "image.jpg"
    }
    return s
}

func UploadGalleryWork(pool *pgxpool.Pool, supabaseURL, supabaseKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		file, header, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no image"})
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "read error"})
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		fileName := fmt.Sprintf("gallery/%s/%d_%s",
		    userID,
		    time.Now().Unix(),
		    sanitizeFileName(header.Filename),
		)

		uploadURL := supabaseURL + "/storage/v1/object/gallery-works/" + fileName

		req, _ := http.NewRequestWithContext(
			context.Background(), "POST", uploadURL,
			bytes.NewReader(data),
		)
		req.Header.Set("Authorization", "Bearer "+supabaseKey)
		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload request failed: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("[gallery] supabase upload error %d: %s\n", resp.StatusCode, string(body))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "supabase upload failed: " + string(body)})
			return
		}

		imageURL := supabaseURL + "/storage/v1/object/public/gallery-works/" + fileName

		title    := c.PostForm("title")
		masterID := c.PostForm("master_id")
		if masterID == "" {
			masterID = userID
		}
		serviceID := c.PostForm("service_id")

		var id string
		pool.QueryRow(c.Request.Context(), `
			INSERT INTO public.gallery_works
			  (master_id, service_id, image_url, title)
			VALUES (
				nullif($1,'')::uuid,
				nullif($2,'')::uuid,
				$3, $4
			)
			RETURNING id::text
		`, masterID, serviceID, imageURL, title).Scan(&id)

		c.JSON(http.StatusCreated, gin.H{"id": id, "image_url": imageURL})
	}
}

func DeleteGalleryWork(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.Parse(c.Param("id"))
		pool.Exec(c.Request.Context(),
			`delete from public.gallery_works where id = $1`, id)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}