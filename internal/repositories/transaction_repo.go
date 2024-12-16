package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/dto"
)

type TransactionRepository interface {
	Transfer(ctx context.Context, from, dest int, userID uuid.UUID) (*dto.AccountTransferResult, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Transfer(ctx context.Context, from, dest int, userID uuid.UUID) (*dto.AccountTransferResult, error) {
	return nil, nil
}
