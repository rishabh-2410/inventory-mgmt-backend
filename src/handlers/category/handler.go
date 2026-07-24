package handlers

import (
	"encoding/json"
	enums "inventory-backend/src/models/enums"
	reqmodels "inventory-backend/src/models/req_dto"
	"inventory-backend/src/requestcontext"
	category "inventory-backend/src/services/category"
	"inventory-backend/src/validation"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	catService *category.Service
}

func New(catService *category.Service) *Handler {
	return &Handler{
		catService: catService,
	}
}

func (h *Handler) HandleAddCategory(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.AddCategoryRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("edit expense request decode failed path=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("edit expense request validation failed path=%s expense_id=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("missing user id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to update category", http.StatusUnauthorized)
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
		slog.Error("cannot add category, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User does not have required permissions to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.catService.AddNewCategory(businessID, &request)
	if err != nil {
		slog.Error("add category failed", "endpoint", r.URL.Path, "user_id", userID, "error", err)
		http.Error(w, "failed to add category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to add category", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleEditCategory(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.UpdateCategoryRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("edit expense request decode failed path=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("edit expense request validation failed path=%s expense_id=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	categoryID := r.PathValue("id")
	if categoryID == "" {
		slog.Error("missing category id", "endpoint", r.URL.Path)
		http.Error(w, "Category ID is required", http.StatusBadRequest)
	}

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("missing user id", "endpoint", r.URL.Path, "category_id", categoryID)
		http.Error(w, "User is not authorized to update category", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("missing business id", "endpoint", r.URL.Path, "category_id", categoryID)
		http.Error(w, "User is either not associated to a business or does not own one", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("invalid user ", "endpoint", r.URL.Path, "category_id", categoryID)
		http.Error(w, "User is either not associated to a business or does not own one", http.StatusUnauthorized)
		return
	}

	if role != enums.RoleOwner {
		slog.Error("cannot update category, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User does not have required permissions to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.catService.UpdateCategory(uuid.MustParse(categoryID), businessID, &request)
	if err != nil {
		slog.Error("edit expense failed path=%s user_id=%s expense_id=%s err=%v", r.URL.Path, userID, categoryID, err)
		http.Error(w, "Unable to update expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to update category", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
