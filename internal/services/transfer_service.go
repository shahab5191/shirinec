package services

import (
	"context"

	"github.com/google/uuid"
	"shirinec.com/internal/repositories"
)

type TransferService interface {
	InternalTransfer(ctx context.Context, from, to int, userID uuid.UUID) (float64, error)
    ExternalTransfer(ctx context.Context, from, ton int, userID uuid.UUID) (float64, error)
}

type transferService struct {
	transferRepo repositories.TransferRepository
}

func NewTransferService(transferRepo repositories.TransferRepository) TransferService {
	return &transferService{
		transferRepo: transferRepo,
	}
}

func (s *transferService) InternalTransfer(ctx context.Context, from, to int, userID uuid.UUID) (float64, error) {
    return 0, nil
}

func (s *transferService) ExternalTransfer(ctx context.Context, from, to int, userID uuid.UUID) (float64, error) {
    return 0, nil
}
