package handlers

import (
	"encoding/json"
	reqmodels "inventory-backend/src/models/req_dto"
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

func (userHandler *Handler) HandleRegisterUser(w http.ResponseWriter,
	r *http.Request) {

	var request reqmodels.RegisterUserRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		slog.Error("failed to decode request body", "path", r.URL.Path, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// logging.Request(r.Method, r.URL.Path, map[string]any{
	// 	"name":     request.Name,
	// 	"email":    request.Email,
	// 	"password": logging.Redacted("password"),
	// })

	err := validation.Validate.Struct(request)
	if err != nil {
		slog.Error("register request validation failed", "path", r.URL.Path, "email", request.Email, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = userHandler.userService.RegisterUser(&request)
	if err != nil {
		if err.Error() == "Email already exists" {
			slog.Error("user is already registered", "path", r.URL.Path, "email", request.Email, "err", err)
			http.Error(w, "User is already registered", http.StatusConflict)
		} else {
			slog.Error("registration failed", "path", r.URL.Path, "email", request.Email, "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// func (userHandler *Handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
// 	var request reqmodels.LoginUserRequestDTO

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()
// 	if err := decoder.Decode(&request); err != nil {
// 		logging.Error("login request decode failed path=%s err=%v", r.URL.Path, err)
// 		http.Error(w, "Invalid Request", http.StatusBadRequest)
// 		return
// 	}

// 	logging.Request(r.Method, r.URL.Path, map[string]any{
// 		"email":    request.Email,
// 		"password": logging.Redacted("password"),
// 	})

// 	err := validation.Validate.Struct(request)
// 	if err != nil {
// 		logging.Error("login request validation failed path=%s email=%s err=%v", r.URL.Path, request.Email, err)
// 		http.Error(w, "Invalid Request", http.StatusBadRequest)
// 		return
// 	}

// 	loggedUser, err := userHandler.userService.LoginUser(&request)
// 	if err != nil {
// 		logging.Error("login request failed path=%s email=%s err=%v", r.URL.Path, request.Email, err)
// 		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	logging.Response(r.Method, r.URL.Path, http.StatusOK, map[string]any{
// 		"id":         loggedUser.ID,
// 		"name":       loggedUser.Name,
// 		"email":      loggedUser.Email,
// 		"created_at": loggedUser.CreatedAt,
// 	})

// 	if err = json.NewEncoder(w).Encode(loggedUser); err != nil {
// 		logging.Error("login response encode failed path=%s email=%s err=%v", r.URL.Path, request.Email, err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }
