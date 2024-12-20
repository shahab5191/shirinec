package services

import (
	"context"

	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/repositories"
)

type IncomeService interface {
	Create(ctx context.Context, input *dto.IncomeCreateRequest, userID uuid.UUID) (*dto.IncomeJoinedResponse, error)
}

type incomeService struct {
	incomeRepo repositories.TransactionRepository
}

func NewIncomeService(incomeRepo repositories.TransactionRepository) IncomeService {
	return &incomeService{
		incomeRepo: incomeRepo,
	}
}

func (s *incomeService) Create(ctx context.Context, input *dto.IncomeCreateRequest, userID uuid.UUID) (*dto.IncomeJoinedResponse, error) {
	return nil, nil
}
