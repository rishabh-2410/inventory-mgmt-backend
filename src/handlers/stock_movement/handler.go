package handlers

import (
	"encoding/json"
	reqmodels "inventory-backend/src/models/req_dto"
	"inventory-backend/src/requestcontext"
	stockmovementservice "inventory-backend/src/services/stock_movement"
	"inventory-backend/src/validation"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	stockMovementService *stockmovementservice.Service
}

func New(stockMovementService *stockmovementservice.Service) *Handler {
	return &Handler{stockMovementService: stockMovementService}
}

func (h *Handler) HandleUpdateStockMovement(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.StockMovementRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("update stock movement request decode failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("update stock movement request validation failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("invalid user: missing user_id", "endpoint", r.URL.Path, "user_id", userID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("invalid user: missing business_id", "endpoint", r.URL.Path, "business_id", businessID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("invalid user: missing user role", "endpoint", r.URL.Path, "role", role)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.stockMovementService.UpdateStockMovement(uuid.MustParse(userID), businessID, &request)
	if err != nil {
		// print err first and then rest of details
		slog.Error("update stock movement failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Unable to update stock movement", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to 	update warehouse", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
