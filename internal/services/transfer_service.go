package services

import (
	"context"

	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/repositories"
)

type TransferService interface {
	Transfer(ctx context.Context, input *dto.TransferRequest, userID uuid.UUID) (*dto.AccountTransferResult, error)
}

type transferService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransferService(transferRepo repositories.TransactionRepository) TransferService {
	return &transferService{
		transactionRepo: transferRepo,
	}
}

func (s *transferService) Transfer(ctx context.Context, input *dto.TransferRequest, userID uuid.UUID) (*dto.AccountTransferResult, error) {
	return nil, nil
}
