package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type IncomeCategoryRepository interface {
	Create(ctx context.Context, category *models.IncomeCategory) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*models.IncomeCategory, error)
	List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.IncomeCategory, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
	Update(ctx context.Context, id int, category *models.IncomeCategory, userID uuid.UUID) error
}

type incomeCategoryRepository struct {
	db *pgxpool.Pool
}

func NewIncomeCategoryRepository(db *pgxpool.Pool) IncomeCategoryRepository {
	return &incomeCategoryRepository{db: db}
}

func (r *incomeCategoryRepository) Create(ctx context.Context, category *models.IncomeCategory) error {
	query := "INSERT INTO income_categories (user_id, name, color) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRow(ctx, query, category.UserID.String(), category.Name, category.Color).Scan(&category.ID)
	return err
}

func (r *incomeCategoryRepository) GetByID(ctx context.Context, ID int, userID uuid.UUID) (*models.IncomeCategory, error) {
	var incomeCategory models.IncomeCategory
	query := "SELECT id, user_id, name, color FROM income_categories WHERE user_id = $1 AND id = $2"
	err := r.db.QueryRow(ctx, query, userID.String(), ID).Scan(&incomeCategory.ID, &incomeCategory.UserID, &incomeCategory.Name, &incomeCategory.Color)
	return &incomeCategory, err
}

func (r *incomeCategoryRepository) List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.IncomeCategory, error) {
	var categories = make([]models.IncomeCategory, 0, limit)
	query := "SELECT id, user_id, name, color FROM income_categories WHERE user_id = $1 LIMIT $2 OFFSET $3"
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return &categories, err
	}
	for rows.Next() {
		var category models.IncomeCategory
		err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Color)
        if err != nil {
            return &categories, err
        }
        categories = append(categories, category)
	}
	return &categories, nil
}

func (r *incomeCategoryRepository) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	return nil
}

func (r *incomeCategoryRepository) Update(ctx context.Context, id int, category *models.IncomeCategory, userID uuid.UUID) error {
	return nil
}
