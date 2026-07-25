package services

import (
	"fmt"
	reqmodels "inventory-backend/src/models/req_dto"
	resmodels "inventory-backend/src/models/res_dto"
	repository "inventory-backend/src/repository/product"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	productRepo *repository.ProductRepo
}

func New(productRepo *repository.ProductRepo) *Service {
	return &Service{
		productRepo: productRepo,
	}
}

func (s *Service) AddNewProduct(businessID uuid.UUID, newProductReq *reqmodels.AddProductRequestDTO) (*resmodels.AddProductResponseDTO, error) {
	result, err := s.productRepo.AddProduct(businessID, newProductReq)
	if err != nil {
		return nil, err
	}

	newProduct := &resmodels.AddProductResponseDTO{
		ID:           result.ID,
		Name:         result.Name,
		BusinessID:   result.BusinessID,
		SKU:          result.SKU,
		SellingPrice: result.SellingPrice,
		CostPrice:    result.CostPrice,
		CategoryID:   result.CategoryID,
		ImageURL:     result.ImageURL,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return newProduct, nil
}

func (s *Service) GetProducts(businessID uuid.UUID) ([]resmodels.AddProductResponseDTO, error) {
	fetchedProducts, err := s.productRepo.GetProducts(businessID)
	if err != nil {
		return nil, err
	}

	result := make([]resmodels.AddProductResponseDTO, 0, len(fetchedProducts))

	for _, item := range fetchedProducts {
		product := resmodels.AddProductResponseDTO{
			ID:           item.ID,
			Name:         item.Name,
			BusinessID:   item.BusinessID,
			SKU:          item.SKU,
			SellingPrice: item.SellingPrice,
			CostPrice:    item.CostPrice,
			CategoryID:   item.CategoryID,
			ImageURL:     item.ImageURL,
		}

		result = append(result, product)
	}

	return result, nil

}

func (s *Service) UpdateProduct(productID uuid.UUID, businessID uuid.UUID, updateProductReq *reqmodels.EditProductRequestDTO) (*resmodels.AddProductResponseDTO, error) {
	if updateProductReq.CategoryID != nil {
		exists, err := s.productRepo.ValidateCategory(*updateProductReq.CategoryID, businessID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, fmt.Errorf("category not found")
		}
	}
	result, err := s.productRepo.UpdateProduct(productID, businessID, updateProductReq)
	if err != nil {
		return nil, err
	}

	newProduct := &resmodels.AddProductResponseDTO{
		ID:           result.ID,
		Name:         result.Name,
		BusinessID:   result.BusinessID,
		SKU:          result.SKU,
		SellingPrice: result.SellingPrice,
		CostPrice:    result.CostPrice,
		CategoryID:   result.CategoryID,
		ImageURL:     result.ImageURL,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return newProduct, nil
}

func (s *Service) DeleteProduct(productID uuid.UUID, businessID uuid.UUID) error {
	slog.Debug("deleting product", "product_id", productID, "business_id", businessID)
	err := s.productRepo.DeleteProduct(productID, businessID)
	if err != nil {
		return err
	}

	slog.Info("product deleted successfully")

	return nil
}
