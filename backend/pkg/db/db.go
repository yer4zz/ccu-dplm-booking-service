package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(databaseURL string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("cannot parse database URL: %v", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
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