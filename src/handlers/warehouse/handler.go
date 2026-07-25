package handlers

import (
	"encoding/json"
	enums "inventory-backend/src/models/enums"
	reqmodels "inventory-backend/src/models/req_dto"
	"inventory-backend/src/requestcontext"
	warehouse "inventory-backend/src/services/warehouse"
	"inventory-backend/src/validation"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	warehouseService *warehouse.Service
}

func New(warehouseService *warehouse.Service) *Handler {
	return &Handler{
		warehouseService: warehouseService,
	}
}

func (h *Handler) HandleAddWarehouse(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.AddWarehouseRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("edit warehouse request decode failed path=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("edit warehouse request validation failed path=%s expense_id=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("missing user id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to update warehouse", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("missing business id", "endpoint", r.URL.Path)
		http.Error(w, "User is either not associated to a business or does not own one", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("invalid user ", "endpoint", r.URL.Path)
		http.Error(w, "User is either not associated to a business or does not own one", http.StatusUnauthorized)
		return
	}

	if role != enums.RoleOwner {
		slog.Error("cannot add warehouse, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User does not have required permissions to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.warehouseService.AddNewWarehouse(businessID, &request)
	if err != nil {
		slog.Error("add warehouse failed", "endpoint", r.URL.Path, "user_id", userID, "error", err)
		http.Error(w, "failed to add warehouse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to add warehouse", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleEditWarehouse(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.EditWarehouseRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("edit warehouse request decode failed path=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("edit warehouse request validation failed path=%s expense_id=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	warehouseID := r.PathValue("id")
	if warehouseID == "" {
		slog.Error("invalid warehouse: missing warehouse id", "endpoint", r.URL.Path)
		http.Error(w, "Invalid Warehouse", http.StatusBadRequest)
	}

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("invalid user: missing user_id", "endpoint", r.URL.Path, "category_id", warehouseID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("invalid user: missing business_id", "endpoint", r.URL.Path, "category_id", warehouseID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("invalid user: missing user role", "endpoint", r.URL.Path, "category_id", warehouseID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	if role != enums.RoleOwner {
		slog.Error("cannot update category, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.warehouseService.UpdateWarehouse(uuid.MustParse(warehouseID), businessID, &request)
	if err != nil {
		// print err first and then rest of details
		slog.Error("edit warehouse failed path=%s user_id=%s expense_id=%s err=%v", r.URL.Path, userID, warehouseID, err)
		http.Error(w, "Unable to update warehouse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to update warehouse", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGetWarehouses(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("invalid user: missing user id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("invalid user: missing business id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	_, ok = requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("invalid user: missing user role ", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.warehouseService.GetWarehouses(businessID)
	if err != nil {
		slog.Error("get warehouses failed", "endpoint", r.URL.Path, "user_id", userID, "error", err)
		http.Error(w, "failed to get warehouses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to get warehouses", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "cannot fetch warehouses", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleDeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("[DeleteWarehouse] invalid user: missing user_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("[DeleteWarehouse]  invalid user: missing business_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("[DeleteWarehouse] invalid user: missing user role", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	if role != enums.RoleOwner {
		slog.Error("[DeleteWarehouse] failed to delete warehouse, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	warehouseID := r.PathValue("id")
	if warehouseID == "" {
		slog.Error("[DeleteWarehouse] missing warehouse id", "endpoint", r.URL.Path)
		http.Error(w, "Invalid item", http.StatusBadRequest)
	}

	err := h.warehouseService.DeleteWarehouse(uuid.MustParse(warehouseID), businessID)
	if err != nil {
		slog.Error("delete warehouse failed", "endpoint", r.URL.Path, "user_id", userID, "error", err)
		http.Error(w, "failed to delete warehouse", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
