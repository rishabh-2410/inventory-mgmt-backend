package handlers

import (
	"encoding/json"
	"errors"
	reqmodels "inventory-backend/src/models/req_dto"
	"inventory-backend/src/requestcontext"
	stockmovementservice "inventory-backend/src/services/stock_movement"
	"inventory-backend/src/validation"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) HandleGetStockMovementHistory(w http.ResponseWriter, r *http.Request) {

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

	requestWithParams, err := h.checkOptionalQueryParams(r)
	if err != nil {
		slog.Error("invalid query params", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid query params", http.StatusBadRequest)
		return
	}

	response, err := h.stockMovementService.GetStockMovementHistory(businessID, requestWithParams)
	if err != nil {
		slog.Error("get stock movement history failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Unable to get stock movement history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to get stock movement history", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) checkOptionalQueryParams(r *http.Request) (*reqmodels.GetStockMovementHistoryRequestDTO, error) {
	var request = &reqmodels.GetStockMovementHistoryRequestDTO{}
	queryParams := r.URL.Query()

	limit := queryParams.Get("limit")
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errors.New("invalid limit")
		}
		request.Limit = &limitInt
	}

	page := queryParams.Get("page")
	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return nil, errors.New("invalid page")
		}
		request.Page = &pageInt
	}

	movementType := queryParams.Get("movement_type")
	if movementType != "" {
		request.MovementType = &movementType
	}

	warehouseID := queryParams.Get("warehouse_id")
	if warehouseID != "" {
		warehouseIDUUID, err := uuid.Parse(warehouseID)
		if err != nil {
			return nil, errors.New("invalid warehouse ID")
		}
		request.WarehouseID = &warehouseIDUUID
	}

	productID := queryParams.Get("product_id")
	if productID != "" {
		productIDUUID, err := uuid.Parse(productID)
		if err != nil {
			return nil, errors.New("invalid product ID")
		}
		request.ProductID = &productIDUUID
	}

	search := queryParams.Get("search")
	if search != "" {
		request.Search = &search
	}

	from := queryParams.Get("from")
	if from != "" {
		fromTime, err := time.Parse(time.DateOnly, from)
		if err != nil {
			return nil, errors.New("invalid from time")
		}
		request.From = &fromTime
	}

	to := queryParams.Get("to")
	if to != "" {
		toTime, err := time.Parse(time.DateOnly, to)
		if err != nil {
			return nil, errors.New("invalid to time")
		}
		request.To = &toTime
	}

	sortBy := queryParams.Get("sort_by")
	if sortBy != "" {
		request.SortBy = &sortBy
	}

	return request, nil
}
