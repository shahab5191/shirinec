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
	CreateForEntity(ctx context.Context, entityTableName string, entityColumn string, media *models.Media, itemID int) error
	CreateForProfile(ctx context.Context, media *models.Media) error
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

func (r *mediaRepository) CreateForEntity(ctx context.Context, entityTableName string, entityColumn string, media *models.Media, itemID int) error {
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

	updateEntityQueryFormat := "UPDATE %s SET %s = $1 WHERE id = $2 AND user_id = $3 RETURNING id"
	updateEntityQuery := fmt.Sprintf(updateEntityQueryFormat, entityTableName, entityColumn)
	log.Println(updateEntityQuery)
	commandTag, err := tx.Exec(ctx, updateEntityQuery, &media.ID, &itemID, &media.UserID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return sql.ErrNoRows
	}

	err = tx.Commit(ctx)
	return err
}

func (r *mediaRepository) CreateForProfile(ctx context.Context, media *models.Media) error {
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

	updateProfileQuery := `
        UPDATE profiles
        SET picture_id = $1
        WHERE id = (SELECT profile_id FROM users WHERE id = $2 AND profile_id IS NOT NULL)
        RETURNING id
    `

	var profileID int
	err = tx.QueryRow(ctx, updateProfileQuery, &media.ID, &media.UserID).Scan(&profileID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}
