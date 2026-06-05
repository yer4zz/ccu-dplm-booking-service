package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(databaseURL string) *pgxpool.Pool {
    pool, err := pgxpool.New(context.Background(), databaseURL)
    if err != nil {
        log.Fatalf("cannot connect to database: %v", err)
    }

    if err := pool.Ping(context.Background()); err != nil {
        log.Fatalf("database ping failed: %v", err)
    }

    log.Println("connected to database")
    return pool
}

func StartKeepalive(pool *pgxpool.Pool) {
    go func() {
        ticker := time.NewTicker(4 * time.Minute)
        for range ticker.C {
            pool.QueryRow(context.Background(), "SELECT 1")
            fmt.Println("[db] keepalive ping")
        }
    }()
}