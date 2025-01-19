package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/src/internal/dto"
	"shirinec.com/src/internal/utils"
)

type TransactionRepository interface {
	Transfer(ctx context.Context, from, dest int, amount float64, userID uuid.UUID) (*dto.AccountTransferResult, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Transfer(ctx context.Context, from, dest int, amount float64, userID uuid.UUID) (*dto.AccountTransferResult, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
    defer func() {
        if err != nil {
            if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
                utils.Logger.Errorf("failed to rollback transaction: %s", rollbackErr.Error())
            }
        }
    }()

	currentTime := time.Now().UTC().Truncate(time.Second)
	createTransactionQuery := `
        INSERT INTO transactions
        (user_id, account_id, amount, transaction_type, update_date, creation_date)
        VALUES($1, $2, $3, $4, $5, $5)
        RETURNING id
    `

	var firstTransID, secondTransID int

	if err := tx.QueryRow(
		ctx,
		createTransactionQuery,
		userID,
		from,
		amount*-1,
		"transfer",
		currentTime,
	).Scan(&firstTransID); err != nil {
		return nil, err
	}

	if err := tx.QueryRow(
		ctx,
		createTransactionQuery,
		userID,
		dest,
		amount,
		"transfer",
		currentTime,
	).Scan(&secondTransID); err != nil {
		return nil, err
	}

	updateTransQuery := `
        UPDATE transactions
        SET linked_transaction_id = $1
        WHERE id = $2
        RETURNING id
    `

	if _, err := tx.Exec(ctx, updateTransQuery, firstTransID, secondTransID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, updateTransQuery, secondTransID, firstTransID); err != nil {
		return nil, err
	}

	changeBalanceQuery := `
        UPDATE accounts
        SET balance = balance + $1
        WHERE id = $2
        RETURNING id, name, type, balance
    `

	var firstAccount, secondAccount dto.AccountTransferResultItem

	if err := tx.QueryRow(ctx, changeBalanceQuery, amount*-1, from).Scan(
		&firstAccount.ID,
		&firstAccount.Name,
		&firstAccount.Type,
		&firstAccount.Balance,
	); err != nil {
		return nil, err
	}
	firstAccount.Change = amount * -1

	if err := tx.QueryRow(ctx, changeBalanceQuery, amount, dest).Scan(
		&secondAccount.ID,
		&secondAccount.Name,
		&secondAccount.Type,
		&secondAccount.Balance,
	); err != nil {
		return nil, err
	}
	secondAccount.Change = amount

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	var result dto.AccountTransferResult
	result.From = firstAccount
	result.Dest = secondAccount
	result.Date = currentTime

	return &result, nil
}
