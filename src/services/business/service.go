package services

import (
	"errors"
	reqmodel "inventory-backend/src/models/req_dto"
	resmodel "inventory-backend/src/models/res_dto"
	repository "inventory-backend/src/repository/business"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	businessRepo *repository.BusinessRepo
}

func New(businessRepo *repository.BusinessRepo) *Service {
	return &Service{
		businessRepo: businessRepo,
	}
}

func (bs *Service) AddBusiness(userID uuid.UUID, createBusinessRequest *reqmodel.CreateBusinessRequestDTO) (*resmodel.CreateBusinessResponseDTO, error) {
	business, err := bs.businessRepo.FindBusinessByName(createBusinessRequest.Name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.Error("failed to find business in db", "name", createBusinessRequest.Name, "err", err)
		return nil, err
	}

	if business != nil {
		slog.Error("failed to add business, duplicate business", "name", createBusinessRequest.Name)
		return nil, errors.New("business name is already taken")

	}

	result, err := bs.businessRepo.CreateBusiness(userID, createBusinessRequest)
	if err != nil {
		slog.Error("failed to add business", "name", createBusinessRequest.Name, "err", err)
		return nil, err
	}

	return &resmodel.CreateBusinessResponseDTO{
		ID:        result.ID,
		Name:      result.Name,
		Email:     result.Email,
		Address:   result.Address,
		Phone:     result.Phone,
		CreatedAt: result.CreatedAt,
	}, nil

}
