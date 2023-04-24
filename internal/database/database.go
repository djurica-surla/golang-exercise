package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Represents connection with the database.
type Connection *sql.DB

// Configuration for creating a new db instance.
type DBConfig struct {
	Host     string
	User     string
	Password string
	DBname   string
}

// Connect connects to the database using the provided DSN.
func Connect(
	ctx context.Context,
	cfg DBConfig,
) (Connection, error) {
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBname)

	instance, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection %w", err)
	}

	err = instance.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database %w", err)
	}

	return instance, nil
}
