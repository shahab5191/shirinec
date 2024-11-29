package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/dto"
	server_errors "shirinec.com/internal/errors"
	"shirinec.com/internal/models"
)

type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.ItemJoinedResponse, error)
	List(ctx context.Context, limit, offset int, userID uuid.UUID) (*[]dto.ItemJoinedResponse, int, error)
    Update(ctx context.Context, item *models.Item) (*dto.ItemJoinedResponse, error)
    Delete(ctx context.Context, id int, userID uuid.UUID) error
}

type itemRepository struct {
	db        *pgxpool.Pool
	tableName string
}

func NewItemRepository(db *pgxpool.Pool) ItemRepository {
	return &itemRepository{db: db, tableName: "items"}
}

func (r *itemRepository) Create(ctx context.Context, item *models.Item) error {
	queryFormat := "INSERT INTO %s (user_id, name, image_id, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	query := fmt.Sprintf(queryFormat, r.tableName)
	// We do not rely on postgres default time because it works with nano seconds but we want second percision
	currentTime := time.Now().UTC().Truncate(time.Second)
	item.CreationDate = currentTime
	item.UpdateDate = currentTime
	err := r.db.QueryRow(ctx, query, item.UserID, item.Name, item.ImageID, item.CategoryID).Scan(&item.ID)
	return err
}

func (r *itemRepository) GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.ItemJoinedResponse, error) {
	queryFormat := "SELECT i.id, i.user_id, i.name, i.image_id, m.url, m.metadata, c.id, c.name, cm.url as category_icon, c.entity_type, i.creation_date, i.update_date FROM %s i LEFT JOIN categories c ON i.category_id = c.id LEFT JOIN media m ON i.image_id = m.id LEFT JOIN media cm ON c.icon_id = m.id WHERE i.id = $1 AND i.user_id = $2"
	query := fmt.Sprintf(queryFormat, r.tableName)

	var item dto.ItemJoinedResponse
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&item.ID, &item.UserID, &item.Name, &item.ImageID, &item.ImageURL, &item.ImageMetadata, &item.CategoryID, &item.CategoryName, &item.CategoryIconURL, &item.CategoryType, &item.CreationDate, &item.UpdateDate)
	return &item, err
}

func (r *itemRepository) List(ctx context.Context, limit, offset int, userID uuid.UUID) (*[]dto.ItemJoinedResponse, int, error) {
	totalCount, err := CountByUserID(ctx, r.db, r.tableName, userID)
	if err != nil {
		return nil, 0, err
	}

	var items = make([]dto.ItemJoinedResponse, 0, limit)
	query := "SELECT i.id, i.user_id, i.name, i.image_id, m.url, m.metadata, c.id, c.name, cm.url as category_icon, c.entity_type, i.creation_date, i.update_date FROM items i LEFT JOIN categories c ON i.category_id = c.id LEFT JOIN media m ON i.image_id = m.id LEFT JOIN media cm ON c.icon_id = m.id WHERE i.user_id = $1 LIMIT $2 OFFSET $3"

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var item dto.ItemJoinedResponse
		err = rows.Scan(&item.ID, &item.UserID, &item.Name, &item.ImageID, &item.ImageURL, &item.ImageMetadata, &item.CategoryID, &item.CategoryName, &item.CategoryIconURL, &item.CategoryType, &item.CreationDate, &item.UpdateDate)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return &items, totalCount, nil
}

func (r *itemRepository) Delete(ctx context.Context, id int, userID uuid.UUID) error {
    queryFormat := "DELETE FROM %s WHERE id = $1 AND user_id = $2 RETURNING id"
    query := fmt.Sprintf(queryFormat, r.tableName)
    var deletedID int
    err := r.db.QueryRow(ctx, query, id, userID).Scan(&deletedID)
    return err
}

func (r *itemRepository) Update(ctx context.Context, item *models.Item) (*dto.ItemJoinedResponse, error) {
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if item.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, item.Name)
		argIndex++
	}

	if item.CategoryID != nil {
		setClauses = append(setClauses, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, item.CategoryID)
		argIndex++
	}

	if item.ImageID != nil {
		setClauses = append(setClauses, fmt.Sprintf("image_id = $%d", argIndex))
		args = append(args, item.ImageID)
		argIndex++
	}

    if len(setClauses) == 0 {
        return nil, &server_errors.EmptyUpdate
    }

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = %d AND user_id = '%s' RETURNING id",
        r.tableName,
		strings.Join(setClauses, ", "),
		item.ID,
		item.UserID.String(),
	)

	err := r.db.QueryRow(ctx, query, args...).Scan(&item.ID)
    if err != nil {
        return nil, err
    }

    itemJoined, err := r.GetByID(ctx, item.ID, item.UserID)

	return itemJoined, err
}
