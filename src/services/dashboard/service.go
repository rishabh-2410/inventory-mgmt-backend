package services

import (
	resmodels "inventory-backend/src/models/res_dto"
	dashboardRepository "inventory-backend/src/repository/dashboard"

	"github.com/google/uuid"
)

type Service struct {
	dashboardRepo *dashboardRepository.DashboardRepo
}

func New(dashboardRepo *dashboardRepository.DashboardRepo) *Service {
	return &Service{dashboardRepo: dashboardRepo}
}

func (s *Service) GetDashboardData(businessID uuid.UUID) (*resmodels.DashboardData, error) {
	return s.dashboardRepo.GetDashboardData(businessID)
}
