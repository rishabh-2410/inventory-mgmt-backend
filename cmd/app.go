package main

import (
	conn "inventory-backend/src/db"
	businessHandler "inventory-backend/src/handlers/business"
	userhandler "inventory-backend/src/handlers/user"
	businessRepo "inventory-backend/src/repository/business"
	tokenRepo "inventory-backend/src/repository/token"
	userRepo "inventory-backend/src/repository/user"
	router "inventory-backend/src/router"
	businessservice "inventory-backend/src/services/business"
	userservice "inventory-backend/src/services/user"
	"inventory-backend/src/validation"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	// Set the global level to Debug immediately
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("starting setup...")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // Local development
	}
	slog.Info("backend configured", "port", port)

	// Initialize validation instance
	validation.InitializeValidator()
	slog.Info("validator initialized")

	db, err := conn.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	slog.Info("database connection ready")

	userRepository := userRepo.NewUserRepo(db)
	tokenRepository := tokenRepo.NewTokenRepo(db)
	businessRepository := businessRepo.NewBusinessRepo(db)

	slog.Info("repositories initialized")

	userservice := userservice.New(userRepository, tokenRepository)
	businessservice := businessservice.New(businessRepository)

	slog.Info("services initialized")

	userHandler := userhandler.New(userservice)
	slog.Info("handlers initialized")

	businessHanlder := businessHandler.New(businessservice)
	slog.Info("handlers initialized")

	router := router.New(userHandler, businessHanlder)

	router.RegisterRoutes()

	slog.Info("listening for incoming requests on", "port", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
