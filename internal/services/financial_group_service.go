package services

import (
	"context"

	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/models"
	"shirinec.com/internal/repositories"
)

type FinancialGroupService interface {
    Create(ctx context.Context, item *dto.FinancialGroupCreateRequest, userID uuid.UUID) (*models.FinancialGroups, error)
}

type financialGroupService struct {
    financialGroupRepo repositories.FinancialGroupRepository
}

func NewFinancialGroupService(financialGroupRepo *repositories.FinancialGroupRepository) FinancialGroupService{
    return &financialGroupService{
        financialGroupRepo: *financialGroupRepo,
    }
}

func (s *financialGroupService) Create(ctx context.Context, item *dto.FinancialGroupCreateRequest, userID uuid.UUID) (*models.FinancialGroups, error) {
    return nil, nil
}

