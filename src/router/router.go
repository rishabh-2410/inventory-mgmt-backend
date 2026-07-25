package router

import (
	"inventory-backend/src/auth"
	business "inventory-backend/src/handlers/business"
	category "inventory-backend/src/handlers/category"
	product "inventory-backend/src/handlers/product"
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
	productHandler   *product.Handler
}

func New(userHandler *user.Handler, businessHandler *business.Handler, categoryHandler *category.Handler, warehouseHandler *warehouse.Handler, productHandler *product.Handler) *Router {
	return &Router{
		userHandler:      userHandler,
		businessHandler:  businessHandler,
		categoryHandler:  categoryHandler,
		warehouseHandler: warehouseHandler,
		productHandler:   productHandler,
	}
}

func (r *Router) RegisterRoutes() {
	slog.Info("registering HTTP routes")
	r.registerAuthRoutes()
	r.registerBusinessRoutes()
	r.registerCategoryRoutes()
	r.registerWarehouseRoutes()
	r.registerProductRoutes()
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

func (r *Router) registerProductRoutes() {
	http.Handle("POST /api/v1/product", auth.Middleware(http.HandlerFunc(r.productHandler.HandleAddProduct)))
	http.Handle("PATCH /api/v1/product/{id}", auth.Middleware(http.HandlerFunc(r.productHandler.HandleEditProduct)))
	http.Handle("GET /api/v1/product", auth.Middleware(http.HandlerFunc(r.productHandler.HandleGetProducts)))
	http.Handle("DELETE /api/v1/product/{id}", auth.Middleware(http.HandlerFunc(r.productHandler.HandleDeleteProduct)))
}
