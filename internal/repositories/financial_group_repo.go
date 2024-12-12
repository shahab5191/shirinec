package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type FinancialGroupRepository interface {
	Create(ctx context.Context, financialGroup *models.FinancialGroups) error
}

type financialGroupRepository struct {
	db *pgxpool.Pool
}

func NewFinancialGroupRepository(db *pgxpool.Pool) FinancialGroupRepository {
	return &financialGroupRepository{
		db: db,
	}
}

func (r *financialGroupRepository) Create(ctx context.Context, financialGroup *models.FinancialGroups) error {
	query := `
        INSERT INTO financial_groups (name, image_id, user_id, creation_date, update_date)
        VALUES ($1, $2, $3, $4, $4)
        RETURNING id
    `
	err := r.db.QueryRow(ctx, query, financialGroup.Name, financialGroup.ImageID, financialGroup.UserID, financialGroup.CreationDate).Scan(&financialGroup.ID)
	return err
}

func (r *financialGroupRepository) AddUserToGroup(ctx context.Context, financialGroupID int, userID uuid.UUID) error {
    query := `
        INSERT INTO user_financial_groups 
        (financial_group_id, user_id)
        VALUES ($1, $2)
    `
    var relID int
    err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&relID)
    return err
}
