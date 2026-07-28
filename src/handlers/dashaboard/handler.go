package handlers

import (
	"encoding/json"
	"inventory-backend/src/requestcontext"
	dashboardService "inventory-backend/src/services/dashboard"
	"log/slog"
	"net/http"
)

type Handler struct {
	dashboardService *dashboardService.Service
}

func New(dashboardService *dashboardService.Service) *Handler {
	return &Handler{dashboardService: dashboardService}
}

func (h *Handler) HandleGetDashboardData(w http.ResponseWriter, r *http.Request) {

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("missing user id", "endpoint", r.URL.Path, "user_id", userID)
		http.Error(w, "User is not authorized to perform this action", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("missing business id", "endpoint", r.URL.Path, "user_id", userID)
		http.Error(w, "User is either not associated to a business or does not own one", http.StatusUnauthorized)
		return
	}

	dashboardData, err := h.dashboardService.GetDashboardData(businessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboardData)
}
