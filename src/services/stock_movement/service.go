package services

import (
	reqmodels "inventory-backend/src/models/req_dto"
	resmodels "inventory-backend/src/models/res_dto"
	repository "inventory-backend/src/repository/stock_movement"

	"github.com/google/uuid"
)

type Service struct {
	stockMovementRepo *repository.StockMovementRepo
}

func New(stockMovementRepo *repository.StockMovementRepo) *Service {
	return &Service{stockMovementRepo: stockMovementRepo}
}

func (s *Service) UpdateStockMovement(userID uuid.UUID, businessID uuid.UUID, stockMovementRequestDTO *reqmodels.StockMovementRequestDTO) (*resmodels.StockMovementResponseDTO, error) {
	currentStock, stockMovement, err := s.stockMovementRepo.UpdateStockMovement(userID, businessID, stockMovementRequestDTO)
	if err != nil {
		return nil, err
	}

	return &resmodels.StockMovementResponseDTO{
		ProductStockID:   stockMovement.ProductStockID,
		ProductID:        stockMovement.ProductID,
		WarehouseID:      stockMovement.WarehouseID,
		CurrentStock:     currentStock,
		MovementType:     stockMovement.MovementType,
		MovementQuantity: stockMovement.MovementQuantity,
		Reason:           stockMovement.Reason,
		CreatedAt:        stockMovement.CreatedAt,
	}, nil
}
