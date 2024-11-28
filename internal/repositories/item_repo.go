package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
    GetByID(ctx context.Context, id int, userID uuid.UUID) (*models.Item, error)
    List(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Item, int, error)
}

type itemRepository struct {
	db *pgxpool.Pool
    tableName string
}

func NewItemRepository(db *pgxpool.Pool) ItemRepository {
	return &itemRepository{db: db, tableName: "items"}
}

func (r *itemRepository) Create(ctx context.Context, item *models.Item) error {
    query := "INSERT INTO $1 (user_id, name, image_id, category_id) VALUES ($2, $3, $4, $5) RETURNING id"
    // We do not rely on postgres default time because it works with nano seconds but we want second percision
    currentTime := time.Now().UTC().Truncate(time.Second)
    item.CreationDate = currentTime
    item.UpdateDate = currentTime
    err := r.db.QueryRow(ctx, query, r.tableName, item.UserID, item.Name, item.ImageID, item.CategoryID).Scan(&item.ID)
    return err
}

func (r *itemRepository) GetByID(ctx context.Context, id int, userID uuid.UUID) (*models.Item, error) {
    query := "SELECT id, user_id, name, image_id, category_id, creation_date, update_date FROM $1 WHERE id = $2 AND user_id = $3"
    var item models.Item
    err := r.db.QueryRow(ctx, query, r.tableName, id, userID).Scan(&item.ID, &item.UserID, &item.Name, &item.ImageID, &item.CategoryID, &item.CreationDate, &item.UpdateDate)
    return &item, err
}

func (r *itemRepository) List(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Item, int, error) {
    totalCount, err := CountByUserID(ctx, r.db, r.tableName, userID)
    if err != nil {
        return nil, 0, err
    }

    var items = make([]models.Item, 0, limit)
    query := "SELECT id, user_id, name, image_id, category_id, creation_date, update_date FROM items WHERE user_id = $1 LIMIT = $2 OFFSET = $3"
    rows, err := r.db.Query(ctx, query, userID, limit, offset)
    defer rows.Close()

    if err != nil {
        return nil, 0, err
    }

    for rows.Next(){
        var item models.Item
        err = rows.Scan(&item.ID, &item.UserID, &item.Name, &item.ImageID, &item.CategoryID, &item.CreationDate, &item.UpdateDate)
        if err != nil {
            return nil, 0, err
        }
        items = append(items, item)
    }
    return &items, totalCount, nil
}
