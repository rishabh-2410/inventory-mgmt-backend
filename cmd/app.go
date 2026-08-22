package main

import (
	conn "inventory-backend/src/db"

	// handlers
	businessHandler "inventory-backend/src/handlers/business"
	categoryhandler "inventory-backend/src/handlers/category"
	dashboardHandler "inventory-backend/src/handlers/dashaboard"
	producthandler "inventory-backend/src/handlers/product"
	productstockhandler "inventory-backend/src/handlers/product_stock"
	stockmovementhandler "inventory-backend/src/handlers/stock_movement"
	userhandler "inventory-backend/src/handlers/user"
	warehouseHandler "inventory-backend/src/handlers/warehouse"

	// repositories
	businessRepo "inventory-backend/src/repository/business"
	categoryRepo "inventory-backend/src/repository/categories"
	dashboardRepo "inventory-backend/src/repository/dashboard"
	productRepo "inventory-backend/src/repository/product"
	productstockRepo "inventory-backend/src/repository/product_stock"
	stockmovementRepo "inventory-backend/src/repository/stock_movement"
	tokenRepo "inventory-backend/src/repository/token"
	userRepo "inventory-backend/src/repository/user"
	warehouseRepo "inventory-backend/src/repository/warehouse"

	// router
	router "inventory-backend/src/router"

	// services
	businessservice "inventory-backend/src/services/business"
	categoryservice "inventory-backend/src/services/category"
	dashboardservice "inventory-backend/src/services/dashboard"
	productservice "inventory-backend/src/services/product"
	productstockservice "inventory-backend/src/services/product_stock"
	stockmovementservice "inventory-backend/src/services/stock_movement"
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

	// Initialize repositories
	userRepository := userRepo.NewUserRepo(db)
	tokenRepository := tokenRepo.NewTokenRepo(db)
	businessRepository := businessRepo.NewBusinessRepo(db)
	categoryRepository := categoryRepo.NewCategoryRepo(db)
	warehouseRepository := warehouseRepo.NewWarehouseRepo(db)
	productRepository := productRepo.NewProductRepo(db)
	stockMovementRepository := stockmovementRepo.NewStockMovementRepo(db)
	dashboardRepository := dashboardRepo.NewDashboardRepo(db)
	productStockRepository := productstockRepo.NewProductStockRepo(db)
	slog.Info("repositories initialized")

	// Initialize services
	userservice := userservice.New(userRepository, tokenRepository)
	businessservice := businessservice.New(businessRepository)
	categoryservice := categoryservice.New(categoryRepository)
	warehouseservice := warehouseservice.New(warehouseRepository)
	productservice := productservice.New(productRepository)
	stockmovementservice := stockmovementservice.New(stockMovementRepository)
	dashboardservice := dashboardservice.New(dashboardRepository)
	productStockService := productstockservice.New(productStockRepository)
	slog.Info("services initialized")

	// Initialize handlers
	userHandler := userhandler.New(userservice)
	slog.Debug("auth handlers initialized")

	businessHanlder := businessHandler.New(businessservice)
	slog.Debug("business handlers initialized")

	categoryHandler := categoryhandler.New(categoryservice)
	slog.Debug("category handlers initialized")

	warehouseHandler := warehouseHandler.New(warehouseservice)
	slog.Debug("warehouse handlers initialized")

	productHandler := producthandler.New(productservice)
	slog.Debug("product handlers initialized")

	stockMovementHandler := stockmovementhandler.New(stockmovementservice)
	slog.Debug("stock movement handlers initialized")

	dashboardHandler := dashboardHandler.New(dashboardservice)
	slog.Debug("dashboard handlers initialized")

	productStockHandler := productstockhandler.New(productStockService)
	slog.Debug("product stock handlers initialized")

	slog.Info("handlers initialized")

	router := router.New(userHandler, businessHanlder, categoryHandler, warehouseHandler, productHandler, stockMovementHandler, dashboardHandler, productStockHandler)

	router.RegisterRoutes()

	slog.Info("listening for incoming requests on", "port", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
