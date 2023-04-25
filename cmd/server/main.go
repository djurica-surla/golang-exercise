package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/djurica-surla/golang-exercise/internal/config"
	"github.com/djurica-surla/golang-exercise/internal/database"
	"github.com/djurica-surla/golang-exercise/internal/service"
	"github.com/djurica-surla/golang-exercise/internal/storage"
	transporthttp "github.com/djurica-surla/golang-exercise/internal/transport/http"
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
	// Second argument is path to migrations file, consider the execution context.
	err = database.Migrate(connection, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate company store.
	companyStore := storage.NewCompanyStore(connection)

	// Instantiate company service.
	companyService := service.NewCompanyService(companyStore)

	// Instantiate mux router.
	router := mux.NewRouter().StrictSlash(true)

	// Instantiate company handler.
	handler := transporthttp.NewCompanyHandler(companyService)

	// Register routes for company handler.
	handler.RegisterRoutes(router)

	// Start the http server.
	log.Printf("starting server on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), router))
}
