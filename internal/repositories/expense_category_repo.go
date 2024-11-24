package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
)

type ExpenseCategoryRepository interface {
	Create(ctx context.Context, category *models.ExpenseCategory) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*models.ExpenseCategory, error)
	List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.ExpenseCategory, int, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
	Update(ctx context.Context, category *models.ExpenseCategory) error
}

type expenseCategoryRepository struct {
	db *pgxpool.Pool
}

func NewExpenseCategoryRepository(db *pgxpool.Pool) ExpenseCategoryRepository {
	return &expenseCategoryRepository{db: db}
}

func (r *expenseCategoryRepository) Create(ctx context.Context, category *models.ExpenseCategory) error {
	query := "INSERT INTO expense_categories (user_id, name, color) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRow(ctx, query, category.UserID.String(), category.Name, category.Color).Scan(&category.ID)
	return err
}

func (r *expenseCategoryRepository) GetByID(ctx context.Context, ID int, userID uuid.UUID) (*models.ExpenseCategory, error) {
	var expenseCategory models.ExpenseCategory
	query := "SELECT id, user_id, name, color FROM expense_categories WHERE user_id = $1 AND id = $2"
	err := r.db.QueryRow(ctx, query, userID.String(), ID).Scan(&expenseCategory.ID, &expenseCategory.UserID, &expenseCategory.Name, &expenseCategory.Color)
	return &expenseCategory, err
}

func (r *expenseCategoryRepository) List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.ExpenseCategory, int, error) {
	var categories = make([]models.ExpenseCategory, 0, limit)
	query := "SELECT id, user_id, name, color FROM expense_categories WHERE user_id = $1 LIMIT $2 OFFSET $3"
	countQuery := "SELECT COUNT(*) FROM expense_categories WHERE user_id = $1"
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
		var category models.ExpenseCategory
		errScan := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Color)
		if errScan != nil {
			return &categories, totalCount, errScan
		}
		categories = append(categories, category)
	}

	if errNext := rows.Err(); errNext != nil {
		log.Printf("[Error] - expenseCategoryRepository.List - iterating rows: %+v\n", errNext)
		return &categories, totalCount, errNext
	}
	return &categories, totalCount, nil
}

func (r *expenseCategoryRepository) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	query := "DELETE FROM expense_categories WHERE id = $1 AND user_id = $2 RETURNING id"
	var deletedID int
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&deletedID)
	return err
}

func (r *expenseCategoryRepository) Update(ctx context.Context, category *models.ExpenseCategory) error {

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
		"UPDATE expense_categories SET %s WHERE id = '%d' AND user_id = '%s'",
		strings.Join(setClauses, ", "),
		category.ID,
		category.UserID,
	)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}
