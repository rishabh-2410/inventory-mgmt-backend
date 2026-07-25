package handlers

import (
	"encoding/json"
	enums "inventory-backend/src/models/enums"
	reqmodels "inventory-backend/src/models/req_dto"
	"inventory-backend/src/requestcontext"
	services "inventory-backend/src/services/product"
	"inventory-backend/src/validation"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	productService *services.Service
}

func New(productService *services.Service) *Handler {
	return &Handler{
		productService: productService,
	}
}

func (h *Handler) HandleAddProduct(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.AddProductRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("edit product request decode failed path=%s err=%v", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("add product request validation failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("[AddProduct] invalid user: missing user_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("[AddProduct] invalid user: missing business_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("[AddProduct] invalid user: missing user role", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	if role != enums.RoleOwner {
		slog.Error("[AddProduct] cannot add product, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User does not have required permissions to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.productService.AddNewProduct(businessID, &request)
	if err != nil {
		slog.Error("[AddProduct] failed to add product", "endpoint", r.URL.Path, "user_id", userID, "error", err)
		http.Error(w, "failed to add product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("[AddProduct] failed to add product", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleEditProduct(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.EditProductRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("[EditProduct] invalid request: decode failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("[EditProduct] invalid request: validation failed", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	productID := r.PathValue("id")
	if productID == "" {
		slog.Error("[EditProduct] invalid request: missing product id", "endpoint", r.URL.Path)
		http.Error(w, "Invalid Warehouse", http.StatusBadRequest)
	}

	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("[EditProduct] invalid user: missing user_id", "endpoint", r.URL.Path, "product_id", productID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("[EditProduct] invalid user: missing business_id", "endpoint", r.URL.Path, "product_id", productID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("[EditProduct] invalid user: missing user role", "endpoint", r.URL.Path, "product_id", productID)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	if role != enums.RoleOwner {
		slog.Error("[EditProduct] cannot update product, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.productService.UpdateProduct(uuid.MustParse(productID), businessID, &request)
	if err != nil {
		// print err first and then rest of details
		slog.Error("[EditProduct] failed to update product", "endpoint", r.URL.Path, "user_id", userID, "product_id", productID, "error", err)
		http.Error(w, "Unable to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("[EditProduct] failed to update product", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("[GetProducts] invalid user: missing user_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("[GetProducts] invalid user: missing business_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	_, ok = requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("[GetProducts] invalid user: missing user role", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	response, err := h.productService.GetProducts(businessID)
	if err != nil {
		slog.Error("[GetProducts] failed to get products", "endpoint", r.URL.Path, "user_id", userID, "error", err)
		http.Error(w, "failed to get products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("[GetProducts] failed to get products", "endpoint", r.URL.Path, "error", err)
		http.Error(w, "cannot fetch products", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("[DeleteProduct] invalid user: missing user_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("[DeleteProduct]  invalid user: missing business_id", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	role, ok := requestcontext.GetUserRole(r.Context())
	if !ok {
		slog.Error("[DeleteProduct] invalid user: missing user role", "endpoint", r.URL.Path)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	if role != enums.RoleOwner {
		slog.Error("[DeleteProduct] failed to delete product, user is not owner", "user_id", userID, "role", role)
		http.Error(w, "User is not authorized to perform this operation", http.StatusUnauthorized)
		return
	}

	productID := r.PathValue("id")
	if productID == "" {
		slog.Error("[DeleteProduct] missing product id", "endpoint", r.URL.Path)
		http.Error(w, "Invalid item", http.StatusBadRequest)
	}

	err := h.productService.DeleteProduct(uuid.MustParse(productID), businessID)
	if err != nil {
		slog.Error("[DeleteProduct] failed to delete product", "endpoint", r.URL.Path, "user_id", userID, "product_id", productID, "error", err)
		http.Error(w, "failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
