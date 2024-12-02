package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/dto"
	"shirinec.com/internal/errors"
	"shirinec.com/internal/models"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error)
	List(ctx context.Context, limit, offset int, userID uuid.UUID) (*[]dto.AccountJoinedResponse, int, error)
	Update(ctx context.Context, account *models.Account) (*dto.AccountJoinedResponse, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
}

type accountRepository struct {
	db        *pgxpool.Pool
	tableName string
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &accountRepository{db: db, tableName: "accounts"}
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
	queryFormat := "INSERT INTO %s (user_id, name, category_id, balance, creation_date, update_date) VALUES ($1, $2, $3, $4, $5, $5) RETURNING id"
	query := fmt.Sprintf(queryFormat, r.tableName)
	currentTime := time.Now().UTC().Truncate(time.Second)
	account.CreationDate = currentTime
	account.UpdateDate = currentTime
	err := r.db.QueryRow(ctx, query, account.UserID, account.Name, account.CategoryID, account.Balance, currentTime).Scan(&account.ID)
	return err
}

func (r *accountRepository) GetByID(ctx context.Context, id int, userID uuid.UUID) (*dto.AccountJoinedResponse, error) {
	queryFormat := "SELECT a.id, a.user_id, a.name, c.id, c.name, c.color, cm.url, a.balance, a.creation_date, a.update_date FROM %s a LEFT JOIN categories c ON a.category_id = c.id LEFT JOIN media cm ON c.icon_id = cm.id WHERE a.id = $1 AND a.user_id = $2"
	query := fmt.Sprintf(queryFormat, r.tableName)

	var item dto.AccountJoinedResponse
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&item.ID, &item.UserID, &item.Name, &item.CategoryID, &item.CategoryName, &item.CategoryColor, &item.CategoryIconURL, &item.Balance, &item.CreationDate, &item.UpdateDate)
	return &item, err
}

func (r *accountRepository) List(ctx context.Context, limit, offset int, userID uuid.UUID) (*[]dto.AccountJoinedResponse, int, error) {
	totalCount, err := CountByUserID(ctx, r.db, r.tableName, userID)
	if err != nil {
		return nil, 0, err
	}
    // TODO
	var accounts = make([]dto.AccountJoinedResponse, 0, limit)
	queryFormat := "SELECT a.id, a.user_id, a.name, c.id, c.name, c.color, cm.url, a.balance, a.creation_date, a.update_date FROM %s a LEFT JOIN categories c ON a.category_id = c.id LEFT JOIN media cm ON c.icon_id = cm.id WHERE a.user_id = $1 LIMIT $2 OFFSET $3"
    query := fmt.Sprintf(queryFormat, r.tableName)

	rows, err := r.db.Query(ctx, query, userID, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.AccountJoinedResponse
		err = rows.Scan(&item.ID, &item.UserID, &item.Name, &item.CategoryID, &item.CategoryName, &item.CategoryColor, &item.CategoryIconURL, &item.Balance, &item.CreationDate, &item.UpdateDate)
		if err != nil {
			return nil, 0, err
		}
		accounts = append(accounts, item)
	}
	return &accounts, totalCount, nil
}

func (r *accountRepository) Delete(ctx context.Context, id int, userID uuid.UUID) error {
	queryFormat := "DELETE FROM %s WHERE id = $1 AND user_id = $2 RETURNING id"
	query := fmt.Sprintf(queryFormat, r.tableName)
	var deletedID int
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&deletedID)
	return err
}

func (r *accountRepository) Update(ctx context.Context, account *models.Account) (*dto.AccountJoinedResponse, error) {
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if account.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, account.Name)
		argIndex++
	}

	if account.CategoryID != nil {
		setClauses = append(setClauses, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, account.CategoryID)
		argIndex++
	}

    if account.Balance != nil {
        setClauses = append(setClauses, fmt.Sprintf("balance = $%d", argIndex))
        args = append(args, account.Balance)
    }

	if len(setClauses) == 0 {
		return nil, &server_errors.EmptyUpdate
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = %d AND user_id = '%s' RETURNING id",
		r.tableName,
		strings.Join(setClauses, ", "),
		account.ID,
		account.UserID.String(),
	)

	err := r.db.QueryRow(ctx, query, args...).Scan(&account.ID)
	if err != nil {
		return nil, err
	}

	accountJoined, err := r.GetByID(ctx, account.ID, account.UserID)

	return accountJoined, err
}
