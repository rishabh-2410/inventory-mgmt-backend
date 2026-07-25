package main

import (
	conn "inventory-backend/src/db"
	businessHandler "inventory-backend/src/handlers/business"
	categoryhandler "inventory-backend/src/handlers/category"
	producthandler "inventory-backend/src/handlers/product"
	userhandler "inventory-backend/src/handlers/user"
	warehouseHandler "inventory-backend/src/handlers/warehouse"
	businessRepo "inventory-backend/src/repository/business"
	categoryRepo "inventory-backend/src/repository/categories"
	productRepo "inventory-backend/src/repository/product"
	tokenRepo "inventory-backend/src/repository/token"
	userRepo "inventory-backend/src/repository/user"
	warehouseRepo "inventory-backend/src/repository/warehouse"
	router "inventory-backend/src/router"
	businessservice "inventory-backend/src/services/business"
	categoryservice "inventory-backend/src/services/category"
	productservice "inventory-backend/src/services/product"
	userservice "inventory-backend/src/services/user"

	warehouseservice "inventory-backend/src/services/warehouse"
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
	categoryRepository := categoryRepo.NewCategoryRepo(db)
	warehouseRepository := warehouseRepo.NewWarehouseRepo(db)
	productRepository := productRepo.NewProductRepo(db)
	slog.Info("repositories initialized")

	userservice := userservice.New(userRepository, tokenRepository)
	businessservice := businessservice.New(businessRepository)
	categoryservice := categoryservice.New(categoryRepository)
	warehouseservice := warehouseservice.New(warehouseRepository)
	productservice := productservice.New(productRepository)
	slog.Info("services initialized")

	userHandler := userhandler.New(userservice)
	slog.Info("handlers initialized")

	businessHanlder := businessHandler.New(businessservice)
	slog.Info("handlers initialized")

	categoryHandler := categoryhandler.New(categoryservice)
	slog.Info("handlers initialized")

	warehouseHandler := warehouseHandler.New(warehouseservice)
	slog.Info("handlers initialized")

	productHandler := producthandler.New(productservice)
	slog.Info("handlers initialized")

	router := router.New(userHandler, businessHanlder, categoryHandler, warehouseHandler, productHandler)

	router.RegisterRoutes()

	slog.Info("listening for incoming requests on", "port", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
