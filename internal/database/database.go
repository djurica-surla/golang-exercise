package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	ErrFailedConnection = errors.New("database connection failed")
	ErrDriver           = errors.New("database migration driver creation failed")
	ErrReadMigration    = errors.New("database migration reading files failed")
	ErrMigration        = errors.New("database migration failed")
)

const (
	PostgresDriver = "postgres"
	// Migrations table name (companies_schema_migrations).
	PostgresMigrationsTable = "companies"
)

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
) (*sql.DB, error) {
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBname)

	db, err := sql.Open(PostgresDriver, postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database %w", err)
	}

	log.Println("Postgres database connection successful")

	return db, nil
}

// Migrate makes sure database migrations are up to date.
func Migrate(connection *sql.DB, path string) error {
	driver, err := postgres.WithInstance(connection,
		&postgres.Config{
			MigrationsTable: fmt.Sprintf("%s_%s", PostgresMigrationsTable, postgres.DefaultMigrationsTable),
		})
	if err != nil {
		return fmt.Errorf("%s: %w", ErrDriver, err)
	}

	// Read migration files from migrations folder in the root.
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", path), PostgresDriver, driver)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrReadMigration, err)
	}

	// Perform database migration.
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: %w", ErrMigration, err)
	} else if err == migrate.ErrNoChange {
		v, _, _ := m.Version()
		log.Printf("postgres migrations up to date, version: %d", v)
	} else if err == nil {
		v, _, _ := m.Version()
		log.Printf("postgres database updated, version: %d", v)
	}

	return nil
}
