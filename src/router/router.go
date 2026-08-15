package router

import (
	"inventory-backend/src/auth"
	business "inventory-backend/src/handlers/business"
	category "inventory-backend/src/handlers/category"
	dashboardHandler "inventory-backend/src/handlers/dashaboard"
	product "inventory-backend/src/handlers/product"
	productstock "inventory-backend/src/handlers/product_stock"
	stockmovement "inventory-backend/src/handlers/stock_movement"
	user "inventory-backend/src/handlers/user"
	warehouse "inventory-backend/src/handlers/warehouse"
	"log/slog"
	"net/http"
)

type Router struct {
	userHandler          *user.Handler
	businessHandler      *business.Handler
	categoryHandler      *category.Handler
	warehouseHandler     *warehouse.Handler
	productHandler       *product.Handler
	stockMovementHandler *stockmovement.Handler
	dashboardHandler     *dashboardHandler.Handler
	productStockHandler  *productstock.Handler
}

func New(
	userHandler *user.Handler,
	businessHandler *business.Handler,
	categoryHandler *category.Handler,
	warehouseHandler *warehouse.Handler,
	productHandler *product.Handler,
	stockMovementHandler *stockmovement.Handler,
	dashboardHandler *dashboardHandler.Handler,
	productStockHandler *productstock.Handler,
) *Router {

	return &Router{
		userHandler:          userHandler,
		businessHandler:      businessHandler,
		categoryHandler:      categoryHandler,
		warehouseHandler:     warehouseHandler,
		productHandler:       productHandler,
		stockMovementHandler: stockMovementHandler,
		dashboardHandler:     dashboardHandler,
		productStockHandler:  productStockHandler,
	}
}

func (r *Router) RegisterRoutes() {
	slog.Info("registering HTTP routes")
	r.registerAuthRoutes()
	r.registerBusinessRoutes()
	r.registerCategoryRoutes()
	r.registerWarehouseRoutes()
	r.registerProductRoutes()
	r.registerStockMovementRoutes()
	r.registerDashboardRoutes()
	r.registerProductStockRoutes()
}

func (r *Router) registerAuthRoutes() {
	http.HandleFunc("POST /api/v1/auth/register", r.userHandler.HandleRegisterUser)
	http.HandleFunc("POST /api/v1/auth/login", r.userHandler.HandleUserLogin)
	http.HandleFunc("POST /api/v1/auth/logout", r.userHandler.HandleUserLogout)
	http.HandleFunc("POST /api/v1/auth/refresh", r.userHandler.HandleTokenRefresh)
	// http.HandleFunc("POST /api/v1/auth/forgot-password", r.userHandler.HandleForgotPassword)
	http.Handle("POST /api/v1/auth/employee", auth.Middleware(http.HandlerFunc(r.userHandler.HandleCreateUser)))
	http.Handle("DELETE /api/v1/auth/employee/{id}", auth.Middleware(http.HandlerFunc(r.userHandler.HandleDeleteUser)))
	http.Handle("GET /api/v1/auth/employee", auth.Middleware(http.HandlerFunc(r.userHandler.HandleGetUsers)))
}

func (r *Router) registerBusinessRoutes() {
	http.HandleFunc("POST /api/v1/business", r.businessHandler.HandleAddBusiness)

}

func (r *Router) registerCategoryRoutes() {
	http.Handle("POST /api/v1/category", auth.Middleware(http.HandlerFunc(r.categoryHandler.HandleAddCategory)))
	http.Handle("PATCH /api/v1/category/{id}", auth.Middleware(http.HandlerFunc(r.categoryHandler.HandleEditCategory)))
	http.Handle("GET /api/v1/category", auth.Middleware(http.HandlerFunc(r.categoryHandler.HandleGetCategories)))
	http.Handle("DELETE /api/v1/category/{id}", auth.Middleware(http.HandlerFunc(r.categoryHandler.HandleDeleteCategory)))

}

func (r *Router) registerWarehouseRoutes() {
	http.Handle("POST /api/v1/warehouse", auth.Middleware(http.HandlerFunc(r.warehouseHandler.HandleAddWarehouse)))
	http.Handle("PATCH /api/v1/warehouse/{id}", auth.Middleware(http.HandlerFunc(r.warehouseHandler.HandleEditWarehouse)))
	http.Handle("GET /api/v1/warehouse", auth.Middleware(http.HandlerFunc(r.warehouseHandler.HandleGetWarehouses)))
	http.Handle("DELETE /api/v1/warehouse/{id}", auth.Middleware(http.HandlerFunc(r.warehouseHandler.HandleDeleteWarehouse)))

}

func (r *Router) registerProductRoutes() {
	http.Handle("POST /api/v1/product", auth.Middleware(http.HandlerFunc(r.productHandler.HandleAddProduct)))
	http.Handle("PATCH /api/v1/product/{id}", auth.Middleware(http.HandlerFunc(r.productHandler.HandleEditProduct)))
	http.Handle("GET /api/v1/product", auth.Middleware(http.HandlerFunc(r.productHandler.HandleGetProducts)))
	http.Handle("DELETE /api/v1/product/{id}", auth.Middleware(http.HandlerFunc(r.productHandler.HandleDeleteProduct)))
	http.Handle("GET /api/v1/product/stock/list", auth.Middleware(http.HandlerFunc(r.productHandler.HandleListProductsWithStock)))
}

func (r *Router) registerStockMovementRoutes() {
	http.Handle("PATCH /api/v1/stock-movement", auth.Middleware(http.HandlerFunc(r.stockMovementHandler.HandleUpdateStockMovement)))
	http.Handle("GET /api/v1/stock-movement", auth.Middleware(http.HandlerFunc(r.stockMovementHandler.HandleGetStockMovementHistory)))
}

func (r *Router) registerDashboardRoutes() {
	http.Handle("GET /api/v1/dashboard", auth.Middleware(http.HandlerFunc(r.dashboardHandler.HandleGetDashboardData)))
	http.Handle("GET /api/v1/inventory/stats", auth.Middleware(http.HandlerFunc(r.dashboardHandler.HandleGetInventoryStats)))
}

func (r *Router) registerProductStockRoutes() {
	http.Handle("GET /api/v1/inventory", auth.Middleware(http.HandlerFunc(r.productStockHandler.HandleGetInventoryList)))
}
