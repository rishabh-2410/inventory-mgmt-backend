package handler

import (
	"encoding/json"
	"errors"
	reqmodels "inventory-backend/src/models/req_dto"
	"inventory-backend/src/requestcontext"
	productstockservice "inventory-backend/src/services/product_stock"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Handler struct {
	productStockService *productstockservice.Service
}

func New(productStockService *productstockservice.Service) *Handler {
	return &Handler{productStockService: productStockService}
}

func (h *Handler) HandleGetInventoryList(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "Invalid selection", http.StatusBadRequest)
		return
	}

	response, err := h.productStockService.GetInventoryList(businessID, requestWithParams)
	if err != nil {
		slog.Error("get inventory list failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Unable to get inventory list", http.StatusInternalServerError)
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

func (h *Handler) checkOptionalQueryParams(r *http.Request) (*reqmodels.GetInventoryListRequestDTO, error) {
	var request = &reqmodels.GetInventoryListRequestDTO{}
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

	return request, nil
}
