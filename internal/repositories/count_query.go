package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CountByUserID(ctx context.Context, db *pgxpool.Pool, tableName string, userID uuid.UUID) (int, error) {
	queryFormat := "SELECT COUNT(*) FROM %s WHERE user_id = $1"
    query := fmt.Sprintf(queryFormat, tableName)
	var totalCount int
	err := db.QueryRow(ctx, query, userID).Scan(&totalCount)
	return totalCount, err
}
