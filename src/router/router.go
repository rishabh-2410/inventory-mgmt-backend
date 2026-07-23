package router

import (
	business "inventory-backend/src/handlers/business"
	user "inventory-backend/src/handlers/user"
	"log/slog"
	"net/http"
)

type Router struct {
	userHandler     *user.Handler
	businessHandler *business.Handler
}

func New(userHandler *user.Handler, businessHandler *business.Handler) *Router {
	return &Router{
		userHandler:     userHandler,
		businessHandler: businessHandler,
	}
}

func (r *Router) RegisterRoutes() {
	slog.Info("registering HTTP routes")
	r.registerAuthRoutes()
	r.registerBusinessRoutes()

}

func (r *Router) registerAuthRoutes() {
	http.HandleFunc("POST /api/v1/auth/register", r.userHandler.HandleRegisterUser)
	http.HandleFunc("POST /api/v1/auth/login", r.userHandler.HandleUserLogin)
}

func (r *Router) registerBusinessRoutes() {
	http.HandleFunc("POST /api/v1/business", r.businessHandler.HandleAddBusiness)
}
