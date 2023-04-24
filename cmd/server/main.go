package main

import (
	"context"
	"log"

	"github.com/djurica-surla/golang-exercise/internal/config"
	"github.com/djurica-surla/golang-exercise/internal/database"
	"github.com/djurica-surla/golang-exercise/internal/storage"
)

func main() {
	// Loads the app config from config.json
	cfg := config.LoadAppConfig()

	// Attempt to establish a connection with the database.
	connection, err := database.Connect(
		context.Background(),
		database.DBConfig{
			Host:     cfg.PostgresConfig.Host,
			User:     cfg.PostgresConfig.User,
			Password: cfg.PostgresConfig.Password,
			DBname:   cfg.PostgresConfig.DBname,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	// Run up migrations to create database schema.
	err = database.Migrate(connection)
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate company store.
	_ = storage.NewCompanyStore(connection)
}
