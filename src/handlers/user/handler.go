package handlers

import (
	"encoding/json"
	reqmodels "inventory-backend/src/models/req_dto"
	"inventory-backend/src/requestcontext"
	user "inventory-backend/src/services/user"
	"inventory-backend/src/validation"
	"log/slog"
	"net/http"
)

type Handler struct {
	userService *user.Service
}

func New(userService *user.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (userHandler *Handler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {

	var request reqmodels.RegisterUserRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		slog.Error("failed to decode request body", "path", r.URL.Path, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err := validation.Validate.Struct(request)
	if err != nil {
		slog.Error("register request validation failed", "path", r.URL.Path, "email", request.Email, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	response, err := userHandler.userService.RegisterUser(&request)
	if err != nil {
		if err.Error() == "email already exists" {
			slog.Error("user is already registered", "path", r.URL.Path, "email", request.Email, "err", err)
			http.Error(w, "User is already registered", http.StatusConflict)
		} else {
			slog.Error("registration failed", "path", r.URL.Path, "email", request.Email, "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("register response encode failed", "path", r.URL.Path, "email", request.Email, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (userHandler *Handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.LoginUserRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		slog.Error("login request decode failed", "path", r.URL.Path, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err := validation.Validate.Struct(request)
	if err != nil {
		slog.Error("login request validation failed", "path", r.URL.Path, "email", request.Email, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	loggedUser, err := userHandler.userService.LoginUser(&request)
	if err != nil {
		slog.Error("login request failed", "path", r.URL.Path, "email", request.Email, "err", err)
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(loggedUser); err != nil {
		slog.Error("login response encode failed", "path", r.URL.Path, "email", request.Email, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (userHandler *Handler) HandleUserLogout(w http.ResponseWriter, r *http.Request) {
	var request *reqmodels.RefreshTokenRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("logout request decode failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err = validation.Validate.Struct(request); err != nil {
		slog.Error("logout request validation failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = userHandler.userService.LogoutUser(request.RefreshToken)
	if err != nil {
		slog.Error("logout request failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (userHandler *Handler) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	var request *reqmodels.RefreshTokenRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("refresh request decode failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err = validation.Validate.Struct(request); err != nil {
		slog.Error("refresh request validation failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	response, err := userHandler.userService.RefreshToken(request.RefreshToken)
	if err != nil {
		slog.Error("refresh request failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("refresh response encode failed path=%s user_id=%s err=%v", r.URL.Path, response.ID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (userHandler *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var request reqmodels.RegisterUserRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("create user request decode failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("create user request validation failed path=%s err=%v", r.URL.Path, err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	_, ok := requestcontext.GetUserID(r.Context())
	if !ok {
		slog.Error("missing user id", "path", r.URL.Path)
		http.Error(w, "User is not authorized to create user", http.StatusUnauthorized)
		return
	}

	businessID, ok := requestcontext.GetBusinessID(r.Context())
	if !ok {
		slog.Error("missing business id", "path", r.URL.Path)
		http.Error(w, "Business is not authorized to create user", http.StatusUnauthorized)
		return
	}

	response, err := userHandler.userService.CreateUser(businessID, &request)
	if err != nil {
		slog.Error("create user request failed", "path", r.URL.Path, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("create user response encode failed", "path", r.URL.Path, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
