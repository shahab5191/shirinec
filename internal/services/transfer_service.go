package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"shirinec.com/internal/dto"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
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
    result, err := s.transactionRepo.Transfer(ctx, input.From, input.Dest, input.Amount, userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows){
            return nil, &server_errors.ItemNotFound
        }

        if pgErr := server_errors.AsPgError(err); pgErr != nil {
            return nil, pgErr
        }

        utils.Logger.Errorf("transferService.Transfer - Calling transactionRepo.Transfer: %s", err.Error())
        return nil, &server_errors.InternalError
    }
	return result, nil
}
