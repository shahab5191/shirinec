package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type FinancialGroupRepository interface {
	Create(ctx context.Context, financialGroup *models.FinancialGroups) error
	AddUserToGroup(ctx context.Context, financialGroupID int, userID uuid.UUID) error
	GetByID(ctx context.Context, finacialGroupID int, userID uuid.UUID) (*models.FinancialGroups, error)
    GetOwnedGroupByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*models.FinancialGroups, error)
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
        RETURNING id
    `
	var relID int
	err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&relID)
	return err
}

func (r *financialGroupRepository) GetByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*models.FinancialGroups, error) {
	var financialGroup models.FinancialGroups
	query := `
        SELECT name, image_id, user_id, creation_date, update_date
        FROM financial_groups
        WHERE id = $1
        AND (
            user_id = $2
            OR
            EXISTS(
                SELECT 1 FROM user_financial_groups
                WHERE financial_group_id = $1
                AND user_id = $2
            )
        )
    `

	err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&financialGroup.Name, &financialGroup.ImageID, &financialGroup.UserID, &financialGroup.CreationDate, &financialGroup.UpdateDate)

	return &financialGroup, err
}

func (r *financialGroupRepository) GetOwnedGroupByID(ctx context.Context, financialGroupID int, userID uuid.UUID) (*models.FinancialGroups, error) {
	var financialGroup models.FinancialGroups
	query := `
        SELECT name, image_id, user_id, creation_date, update_date
        FROM financial_groups
        WHERE id = $1 AND user_id = $2
    `

	err := r.db.QueryRow(ctx, query, financialGroupID, userID).Scan(&financialGroup.Name, &financialGroup.ImageID, &financialGroup.UserID, &financialGroup.CreationDate, &financialGroup.UpdateDate)

	return &financialGroup, err
}
