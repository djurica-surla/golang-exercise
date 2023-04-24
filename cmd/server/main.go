package main

import (
	"context"
	"log"

	"github.com/djurica-surla/golang-exercise/internal/config"
	"github.com/djurica-surla/golang-exercise/internal/database"
)

func main() {
	// Loads the app config from config.json
	cfg := config.LoadAppConfig()

	// Attempt to establish a connection with the database.
	_, err := database.Connect(
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
	log.Println("Postgres database connection successful")
}
