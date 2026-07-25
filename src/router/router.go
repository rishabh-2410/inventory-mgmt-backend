package router

import (
	"inventory-backend/src/auth"
	business "inventory-backend/src/handlers/business"
	category "inventory-backend/src/handlers/category"
	user "inventory-backend/src/handlers/user"
	warehouse "inventory-backend/src/handlers/warehouse"
	"log/slog"
	"net/http"
)

type Router struct {
	userHandler      *user.Handler
	businessHandler  *business.Handler
	categoryHandler  *category.Handler
	warehouseHandler *warehouse.Handler
}

func New(userHandler *user.Handler, businessHandler *business.Handler, categoryHandler *category.Handler, warehouseHandler *warehouse.Handler) *Router {
	return &Router{
		userHandler:      userHandler,
		businessHandler:  businessHandler,
		categoryHandler:  categoryHandler,
		warehouseHandler: warehouseHandler,
	}
}

func (r *Router) RegisterRoutes() {
	slog.Info("registering HTTP routes")
	r.registerAuthRoutes()
	r.registerBusinessRoutes()
	r.registerCategoryRoutes()
	r.registerWarehouseRoutes()

}

func (r *Router) registerAuthRoutes() {
	http.HandleFunc("POST /api/v1/auth/register", r.userHandler.HandleRegisterUser)
	http.HandleFunc("POST /api/v1/auth/login", r.userHandler.HandleUserLogin)
	http.HandleFunc("POST /api/v1/auth/logout", r.userHandler.HandleUserLogout)
	http.HandleFunc("POST /api/v1/auth/refresh", r.userHandler.HandleTokenRefresh)
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
