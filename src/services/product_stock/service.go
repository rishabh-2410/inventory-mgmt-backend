package services

import (
	reqmodels "inventory-backend/src/models/req_dto"
	resmodels "inventory-backend/src/models/res_dto"
	repository "inventory-backend/src/repository/product_stock"

	"github.com/google/uuid"
)

type Service struct {
	productStockRepo *repository.ProductStockRepo
}

func New(productStockRepo *repository.ProductStockRepo) *Service {
	return &Service{productStockRepo: productStockRepo}
}

func (s *Service) GetInventoryList(businessID uuid.UUID, requestWithParams *reqmodels.GetInventoryListRequestDTO) ([]*resmodels.ProductStockResponseDTO, error) {
	err := s.validateRequestWithParams(requestWithParams)
	if err != nil {
		return nil, err
	}

	offset := 0
	if requestWithParams.Page != nil && requestWithParams.Limit != nil {
		offset = max((*requestWithParams.Page-1)*(*requestWithParams.Limit), 0)
	}

	productStockList, err := s.productStockRepo.GetInventoryList(businessID, requestWithParams, offset)
	if err != nil {
		return nil, err
	}

	productStockResponseDTO := make([]*resmodels.ProductStockResponseDTO, 0, len(productStockList))
	for _, productStock := range productStockList {
		productStockResponseDTO = append(productStockResponseDTO, &resmodels.ProductStockResponseDTO{
			ID:            productStock.ID,
			ProductName:   productStock.ProductName,
			Sku:           productStock.Sku,
			CategoryName:  productStock.CategoryName,
			WarehouseName: productStock.WarehouseName,
			CurrentStock:  productStock.CurrentStock,
			UpdatedAt:     productStock.UpdatedAt,
		})
	}

	return productStockResponseDTO, nil
}

func (s *Service) validateRequestWithParams(requestWithParams *reqmodels.GetInventoryListRequestDTO) error {
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

	searchParam, err := ValidateSearchParam(requestWithParams.Search)
	if err != nil {
		return err
	}
	requestWithParams.Search = searchParam

	return nil
}
