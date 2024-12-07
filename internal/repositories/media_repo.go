package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
	"shirinec.com/internal/utils"
)

type MediaRepository interface {
	Create(ctx context.Context, media *models.Media) error
	ListForCleanUp(ctx context.Context, threshold string) ([]string, error)
	DeleteRemovedMedia(ctx context.Context) error
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

func (r *mediaRepository) ListForCleanUp(ctx context.Context, threshold string) ([]string, error) {

	queryList := `
        UPDATE media
        SET status = 'removed'
        WHERE status = 'temp'
        AND creation_date <= NOW() - $1::interval
        RETURNING file_path;
    `

	rows, err := r.db.Query(ctx, queryList, threshold)
	if err != nil {
		return nil, err
	}
    defer rows.Close()

	var list []string
	for rows.Next() {
        utils.Logger.Info("row")
		var fileName string
		err := rows.Scan(&fileName)
		if err != nil {
			return nil, err
		}
		list = append(list, fileName)
	}

    if err := rows.Err(); err != nil {
        return nil, err
    }
	return list, nil
}

func (r *mediaRepository) DeleteRemovedMedia(ctx context.Context) error {
	query := "DELETE FROM media WHERE status = 'removed'"
	if _, err := r.db.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}
