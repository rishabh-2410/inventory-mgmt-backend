package handlers

import (
	"encoding/json"
	reqmodel "inventory-backend/src/models/req_dto"
	businessService "inventory-backend/src/services/business"
	"inventory-backend/src/validation"
	"log/slog"
	"net/http"
)

type Handler struct {
	businessService *businessService.Service
}

func New(service *businessService.Service) *Handler {
	return &Handler{
		businessService: service,
	}
}

func (bh *Handler) HandleAddBusiness(w http.ResponseWriter, r *http.Request) {
	var request *reqmodel.CreateBusinessRequestDTO

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		slog.Error("invalid request", "request", *request, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(request)
	if err != nil {
		slog.Error("invalid request", "request", *request, "err", err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	response, err := bh.businessService.AddBusiness(request)

	if err != nil {
		if err.Error() == "business name is already taken" {
			http.Error(w, "Business name is already taken, try with another name", http.StatusConflict)
		} else {
			slog.Error("failed to register business", "path", r.URL.Path, "name", request.Name, "err", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
