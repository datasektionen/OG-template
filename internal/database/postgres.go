package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

// Init connects to PostgreSQL and applies the schema SQL.
func Init(connStr string, schemaSQL string) (*sql.DB, error) {
	if connStr == "" {
		return nil, nil
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("Successfully connected to database")

	if schemaSQL != "" {
		if _, err := db.ExecContext(ctx, schemaSQL); err != nil {
			return nil, fmt.Errorf("failed to apply schema.sql: %w", err)
		}
		slog.Info("Database schema applied successfully")
	}

	return db, nil
}
