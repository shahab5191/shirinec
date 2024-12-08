package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/models"
)

type MediaRepository interface {
	Create(ctx context.Context, media *models.Media) error
	ListForCleanUp(ctx context.Context, threshold string) ([]string, error)
	DeleteRemovedMedia(ctx context.Context) error
    GetByMediaName(ctx context.Context, url string, userID uuid.UUID) (*models.Media, error)
}

type mediaRepository struct {
	db        *pgxpool.Pool
	tableName string
}

func NewMediaRepository(db *pgxpool.Pool) MediaRepository {
	return &mediaRepository{db: db, tableName: "media"}
}

func (r *mediaRepository) Create(ctx context.Context, media *models.Media) error {
	queryFormat := "INSERT INTO %s (url, file_path, user_id, metadata, creation_date, update_date, access, financial_group_id) VALUES ($1, $2, $3, $4, $5, $5, $6, $7) RETURNING id"
	query := fmt.Sprintf(queryFormat, r.tableName)

	err := r.db.QueryRow(ctx, query, &media.Url, &media.FilePath, &media.UserID, &media.Metadata, &media.CreationDate, &media.Access, &media.FinancialGroupID).Scan(&media.ID)
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

func (r *mediaRepository) GetByMediaName(ctx context.Context, url string, userID uuid.UUID) (*models.Media, error) {
    query := `
        select m.id, m.url, m.file_path, m.metadata, m.access, m.financial_group_id, m.creation_date, m.update_date from media m
        left join user_financial_groups ufg
        on (
            ufg.financial_group_id = m.financial_group_id
            and
            ufg.user_id = $1
        )
        where 
        m.url = $2
        and
        (
            m.user_id = $1
            or
            m.access = 'public'
            or
            (m.access = 'group' and ufg.id is not null)
        )`
    var media models.Media
    err := r.db.QueryRow(ctx, query, userID, url).Scan(&media.ID, &media.Url, &media.FilePath, &media.Metadata, &media.Access, &media.FinancialGroupID, &media.CreationDate, &media.UpdateDate)
    return &media, err
}

