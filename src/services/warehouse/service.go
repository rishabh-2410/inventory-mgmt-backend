package services

import (
	reqmodels "inventory-backend/src/models/req_dto"
	resmodels "inventory-backend/src/models/res_dto"
	repository "inventory-backend/src/repository/warehouse"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	warehouseRepo *repository.WarehouseRepo
}

func New(warehouseRepo *repository.WarehouseRepo) *Service {
	return &Service{
		warehouseRepo: warehouseRepo,
	}
}

func (s *Service) AddNewWarehouse(businessID uuid.UUID, newWarehouseReq *reqmodels.AddWarehouseRequestDTO) (*resmodels.AddWarehouseResponseDTO, error) {
	result, err := s.warehouseRepo.AddWarehouse(businessID, newWarehouseReq)
	if err != nil {
		return nil, err
	}

	newWarehouse := &resmodels.AddWarehouseResponseDTO{
		ID:         result.ID,
		Name:       result.Name,
		BusinessID: result.BusinessID,
		Address:    result.Address,
	}

	return newWarehouse, nil
}

func (s *Service) GetWarehouses(businessID uuid.UUID) ([]resmodels.AddWarehouseResponseDTO, error) {
	fetchedWarehouses, err := s.warehouseRepo.GetWarehouses(businessID)
	if err != nil {
		return nil, err
	}

	result := make([]resmodels.AddWarehouseResponseDTO, 0, len(fetchedWarehouses))

	for _, item := range fetchedWarehouses {
		warehouse := resmodels.AddWarehouseResponseDTO{
			ID:         item.ID,
			Name:       item.Name,
			BusinessID: item.BusinessID,
			Address:    item.Address,
		}

		result = append(result, warehouse)
	}

	return result, nil

}

func (s *Service) UpdateWarehouse(categoryID uuid.UUID, businessID uuid.UUID, updateWarehouseReq *reqmodels.EditWarehouseRequestDTO) (*resmodels.AddWarehouseResponseDTO, error) {
	result, err := s.warehouseRepo.UpdateWarehouse(categoryID, businessID, updateWarehouseReq)
	if err != nil {
		return nil, err
	}

	newWarehouse := &resmodels.AddWarehouseResponseDTO{
		ID:         result.ID,
		Name:       result.Name,
		BusinessID: result.BusinessID,
		Address:    result.Address,
	}

	return newWarehouse, nil
}

func (s *Service) DeleteWarehouse(warehouseID uuid.UUID, businessID uuid.UUID) error {
	slog.Debug("deleting warehouse", "warehouse_id", warehouseID, "business_id", businessID)
	err := s.warehouseRepo.DeleteWarehouse(warehouseID, businessID)
	if err != nil {
		return err
	}

	slog.Info("warehouse deleted successfully")

	return nil
}
