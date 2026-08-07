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

func (s *Service) GetStockMovementHistory(businessID uuid.UUID, requestWithParams *reqmodels.GetStockMovementHistoryRequestDTO) (*resmodels.StockMovementHistoryResponseDTO, error) {
	err := s.validateRequestWithParams(requestWithParams)
	if err != nil {
		return nil, err
	}

	stockMovementHistoryResponseDTOs := &resmodels.StockMovementHistoryResponseDTO{}
	offset := 0
	if requestWithParams.Page != nil && requestWithParams.Limit != nil {
		offset = max((*requestWithParams.Page-1)*(*requestWithParams.Limit), 0)
		stockMovementHistoryResponseDTOs.Page = *requestWithParams.Page
		stockMovementHistoryResponseDTOs.Limit = *requestWithParams.Limit
	}

	stockMovementsHistory, err := s.stockMovementRepo.GetStockMovementHistory(businessID, requestWithParams, offset)
	if err != nil {
		return nil, err
	}

	for _, stockMovementHistory := range stockMovementsHistory {
		stockMovementHistoryResponseDTOs.Data = append(stockMovementHistoryResponseDTOs.Data, &resmodels.StockMovementHistory{
			ID:               stockMovementHistory.ID,
			UserID:           stockMovementHistory.UserID,
			ProductID:        stockMovementHistory.ProductID,
			MovementType:     stockMovementHistory.MovementType,
			MovementQuantity: stockMovementHistory.MovementQuantity,
			Reason:           stockMovementHistory.Reason,
			CreatedAt:        stockMovementHistory.CreatedAt,
			ProductName:      stockMovementHistory.ProductName,
			Sku:              stockMovementHistory.Sku,
			WarehouseName:    stockMovementHistory.WarehouseName,
		})
	}
	stockMovementHistoryResponseDTOs.Total = len(stockMovementsHistory)
	stockMovementHistoryResponseDTOs.HasNext = stockMovementHistoryResponseDTOs.Limit < stockMovementHistoryResponseDTOs.Total

	if stockMovementHistoryResponseDTOs.HasNext {
		stockMovementHistoryResponseDTOs.Data = stockMovementHistoryResponseDTOs.Data[:stockMovementHistoryResponseDTOs.Limit]
	}

	return stockMovementHistoryResponseDTOs, nil
}

func (s *Service) validateRequestWithParams(requestWithParams *reqmodels.GetStockMovementHistoryRequestDTO) error {
	limit, err := ValidateLimit(requestWithParams.Limit)
	if err != nil {
		return err
	}
	requestWithParams.Limit = limit

	page, err := ValidatePage(requestWithParams.Page)
	if err != nil {
		return err
	}
	requestWithParams.Page = page

	sortBy, err := ValidateSortBy(requestWithParams.SortBy)
	if err != nil {
		return err
	}
	requestWithParams.SortBy = sortBy

	from, to, err := ValidateTimeRange(requestWithParams.From, requestWithParams.To)
	if err != nil {
		return err
	}
	requestWithParams.From = from
	requestWithParams.To = to

	movementType, err := ValidateMovementType(requestWithParams.MovementType)
	if err != nil {
		return err
	}
	requestWithParams.MovementType = movementType

	searchParam, err := ValidateSearchParam(requestWithParams.Search)
	if err != nil {
		return err
	}
	requestWithParams.Search = searchParam

	return nil
}
