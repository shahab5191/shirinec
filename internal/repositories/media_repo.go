package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type MediaRepository interface {
	Create(ctx context.Context, media *models.Media) error
	CreateForItem(ctx context.Context, media *models.Media, itemID int) error
}

type mediaRepository struct {
	db        *pgxpool.Pool
	tableName string
}

func NewMediaRepository(db *pgxpool.Pool) MediaRepository {
	return &mediaRepository{db: db, tableName: "media"}
}

func (r *mediaRepository) Create(ctx context.Context, media *models.Media) error {
	queryFormat := "INSERT INTO %s (url, file_path, user_id, metadata, creation_date, update_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	query := fmt.Sprintf(queryFormat, r.tableName)

	err := r.db.QueryRow(ctx, query, &media.Url, &media.FilePath, &media.UserID, &media.Metadata, &media.CreationDate, &media.UpdateDate).Scan(&media.ID)
	if err != nil {
		log.Printf("[Error] - mediaRepository.Create - Running query: %+v\n", err)
	}
	return err
}

func (r *mediaRepository) CreateForItem(ctx context.Context, media *models.Media, itemID int) error {
	queryFormat := "INSERT INTO %s (url, file_path, user_id, metadata, creation_date, update_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	query := fmt.Sprintf(queryFormat, r.tableName)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("[Error] - mediaRepository - Begining Transaction: %+v\n", err)
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, query, &media.Url, &media.FilePath, &media.UserID, &media.Metadata, &media.CreationDate, &media.UpdateDate).Scan(&media.ID)
	if err != nil {
		return err
	}

	updateItemQuery := "UPDATE items SET image_id = $1 WHERE id = $2 AND user_id = $3 RETURNING id"
	commandTag, err := tx.Exec(ctx, updateItemQuery, &media.ID, &itemID, &media.UserID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return sql.ErrNoRows
	}

	err = tx.Commit(ctx)
	return err
}
