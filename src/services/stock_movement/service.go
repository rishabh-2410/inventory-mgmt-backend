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

func (s *Service) GetStockMovementHistory(businessID uuid.UUID) ([]*resmodels.StockMovementHistoryResponseDTO, error) {
	stockMovementsHistory, err := s.stockMovementRepo.GetStockMovementHistory(businessID)
	if err != nil {
		return nil, err
	}

	stockMovementHistoryResponseDTOs := []*resmodels.StockMovementHistoryResponseDTO{}
	for _, stockMovementHistory := range stockMovementsHistory {
		stockMovementHistoryResponseDTOs = append(stockMovementHistoryResponseDTOs, &resmodels.StockMovementHistoryResponseDTO{
			UserID:           stockMovementHistory.UserID,
			WarehouseID:      stockMovementHistory.WarehouseID,
			ProductID:        stockMovementHistory.ProductID,
			MovementType:     stockMovementHistory.MovementType,
			MovementQuantity: stockMovementHistory.MovementQuantity,
			Reason:           stockMovementHistory.Reason,
			CreatedAt:        stockMovementHistory.CreatedAt,
		})
	}

	return stockMovementHistoryResponseDTOs, nil
}
