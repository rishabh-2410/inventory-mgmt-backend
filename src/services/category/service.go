package services

import (
	reqmodels "inventory-backend/src/models/req_dto"
	resmodels "inventory-backend/src/models/res_dto"
	repository "inventory-backend/src/repository/categories"
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	categoryRepo *repository.CategoryRepo
}

func New(categoryRepo *repository.CategoryRepo) *Service {
	return &Service{
		categoryRepo: categoryRepo,
	}
}

func (s *Service) AddNewCategory(businessID uuid.UUID, newCategoryReq *reqmodels.AddCategoryRequestDTO) (*resmodels.AddCategoryResponseDTO, error) {
	result, err := s.categoryRepo.AddCategory(businessID, newCategoryReq)
	if err != nil {
		return nil, err
	}

	newCategory := &resmodels.AddCategoryResponseDTO{
		ID:          result.ID,
		Name:        result.Name,
		BusinessID:  result.BusinessID,
		Description: result.Description,
	}

	return newCategory, nil
}

func (s *Service) GetCategories(businessID uuid.UUID) ([]resmodels.AddCategoryResponseDTO, error) {
	fetchedCategories, err := s.categoryRepo.GetCategories(businessID)
	if err != nil {
		return nil, err
	}

	result := make([]resmodels.AddCategoryResponseDTO, 0, len(fetchedCategories))

	for _, item := range fetchedCategories {
		category := resmodels.AddCategoryResponseDTO{
			ID:          item.ID,
			Name:        item.Name,
			BusinessID:  item.BusinessID,
			Description: item.Description,
		}

		result = append(result, category)
	}

	return result, nil

}

func (s *Service) UpdateCategory(categoryID uuid.UUID, businessID uuid.UUID, updateCategoryReq *reqmodels.UpdateCategoryRequestDTO) (*resmodels.AddCategoryResponseDTO, error) {
	result, err := s.categoryRepo.UpdateCategory(categoryID, businessID, updateCategoryReq)
	if err != nil {
		return nil, err
	}

	newCategory := &resmodels.AddCategoryResponseDTO{
		ID:          result.ID,
		Name:        result.Name,
		BusinessID:  result.BusinessID,
		Description: result.Description,
	}

	return newCategory, nil
}

func (s *Service) DeleteCategory(categoryID uuid.UUID, businessID uuid.UUID) error {
	slog.Debug("deleting category", "category_id", categoryID, "business_id", businessID)
	err := s.categoryRepo.DeleteCategory(categoryID, businessID)
	if err != nil {
		return err
	}

	slog.Info("category deleted successfully")

	return nil
}
