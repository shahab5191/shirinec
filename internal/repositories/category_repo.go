package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/models"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*models.Category, error)
	List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.Category, int, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
	Update(ctx context.Context, category *models.Category) error
}

type categoryRepository struct {
	db *pgxpool.Pool
    tableName string
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db, tableName: "categories"}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	queryFormat := "INSERT INTO %s (user_id, name, color, icon_id, entity_type, update_date, creation_date) VALUES ($1, $2, $3, $4, $5, $6, $6) RETURNING id"
    query := fmt.Sprintf(queryFormat, r.tableName)
	currentTime := time.Now().UTC().Truncate(time.Second)
	category.CreationDate = &currentTime
	category.UpdateDate = &currentTime
	err := r.db.QueryRow(ctx, query, category.UserID.String(), category.Name, category.Color, category.IconID, category.EntityType, currentTime).Scan(&category.ID)
	return err
}

func (r *categoryRepository) GetByID(ctx context.Context, ID int, userID uuid.UUID) (*models.Category, error) {
	var category models.Category
	queryFormat := "SELECT id, user_id, name, color, icon_id, entity_type, creation_date, update_date FROM %s WHERE user_id = $1 AND id = $2"
    query := fmt.Sprintf(queryFormat, r.tableName)

	err := r.db.QueryRow(ctx, query, userID.String(), ID).Scan(&category.ID, &category.UserID, &category.Name, &category.Color, &category.IconID, &category.EntityType, &category.CreationDate, &category.UpdateDate)
	return &category, err
}

func (r *categoryRepository) List(ctx context.Context, limit int, offset int, userID uuid.UUID) (*[]models.Category, int, error) {
    totalCount, err := CountByUserID(ctx, r.db, r.tableName, userID)
    if err != nil {
        return nil, 0, err
    }

	var categories = make([]models.Category, 0, limit)
	queryFormat := "SELECT id, user_id, name, color, icon_id, entity_type, creation_date, update_date FROM %s WHERE user_id = $1 LIMIT $2 OFFSET $3"
    query := fmt.Sprintf(queryFormat, r.tableName)
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	defer rows.Close()

	if err != nil {
		return &categories, totalCount, err
	}

	for rows.Next() {
		var category models.Category
		errScan := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Color, &category.IconID, &category.EntityType, &category.CreationDate, &category.UpdateDate)
		if errScan != nil {
			return &categories, totalCount, errScan
		}
		categories = append(categories, category)
	}

	if errNext := rows.Err(); errNext != nil {
		log.Printf("[Error] - categoryRepository.List - iterating rows: %+v\n", errNext)
		return &categories, totalCount, errNext
	}
	return &categories, totalCount, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	queryFormat := "DELETE FROM %s WHERE id = $1 AND user_id = $2 RETURNING id"
    query := fmt.Sprintf(queryFormat, r.tableName)
	var deletedID int
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&deletedID)
	return err
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {

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

	if category.IconID != nil {
		setClauses = append(setClauses, fmt.Sprintf("icon_id = $%d", argIndex))
		args = append(args, *category.IconID)
		argIndex++
	}

	setClauses = append(setClauses, fmt.Sprintf("update_date = $%d", argIndex))
	args = append(args, time.Now().UTC())
	argIndex++

	if len(setClauses) == 0 {
        log.Printf("No update provided!: %+v\n", setClauses)
		return &server_errors.EmptyUpdate
	}
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = '%d' AND user_id = '%s' RETURNING id",
        r.tableName,
		strings.Join(setClauses, ", "),
		category.ID,
		category.UserID,
	)
    
	err := r.db.QueryRow(ctx, query, args...).Scan(&category.ID)
	return err
}
