package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/models"
)

type IncomeCategoryRepository interface {
	Create(ctx context.Context, category *models.IncomeCategory) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*models.IncomeCategory, error)
	List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.IncomeCategory, int, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
	Update(ctx context.Context, category *models.IncomeCategory) error
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

func (r *incomeCategoryRepository) List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.IncomeCategory, int, error) {
	var categories = make([]models.IncomeCategory, 0, limit)
	query := "SELECT id, user_id, name, color FROM income_categories WHERE user_id = $1 LIMIT $2 OFFSET $3"
	countQuery := "SELECT COUNT(*) FROM income_categories WHERE user_id = $1"
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	defer rows.Close()

	if err != nil {
		return &categories, totalCount, err
	}

	for rows.Next() {
		var category models.IncomeCategory
		errScan := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Color)
		if errScan != nil {
			return &categories, totalCount, errScan
		}
		categories = append(categories, category)
	}

	if errNext := rows.Err(); errNext != nil {
		log.Printf("[Error] - incomeCategoryRepository.List - iterating rows: %+v\n", errNext)
		return &categories, totalCount, errNext
	}
	return &categories, totalCount, nil
}

func (r *incomeCategoryRepository) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	query := "DELETE FROM income_categories WHERE id = $1 AND user_id = $2 RETURNING id"
	var deletedID int
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&deletedID)
	return err
}

func (r *incomeCategoryRepository) Update(ctx context.Context, category *models.IncomeCategory) error {

	var setClauses []string
	var args []interface{}
	argIndex := 1

	if category.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *category.Name)
		argIndex++
	}

	if category.Color != nil {
		setClauses = append(setClauses, fmt.Sprintf("color = $%d", argIndex))
		args = append(args, *category.Color)
		argIndex++
	}

	if len(setClauses) == 0 {
		return &server_errors.EmptyUpdate
	}
	query := fmt.Sprintf(
		"UPDATE income_categories SET %s WHERE id = '%d' AND user_id = '%s'",
		strings.Join(setClauses, ", "),
		category.ID,
		category.UserID,
	)

    _, err := r.db.Exec(ctx, query, args...)
	return err
}
