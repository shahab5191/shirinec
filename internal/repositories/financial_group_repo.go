package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type FinancialGroupRepository interface{
    Create(ctx context.Context, financialGroup *models.FinancialGroups) error
}

type financialGroupRepository struct {
    db *pgxpool.Pool
}

func NewFinancialGroupRepository(db *pgxpool.Pool) FinancialGroupRepository {
    return &financialGroupRepository {
        db: db,
    }
}

func (r *financialGroupRepository) Create(ctx context.Context, financialGroup *models.FinancialGroups) error {
    return nil
}
