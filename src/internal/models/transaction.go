package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                  int
	UserID              uuid.UUID
	AccountID           int
	CategoryID          int
	Amount              float32
	Description         string
	TransactionType     TransactionTypes
	LinkedTransactionID int
	CreationDate        time.Time
	UpdateDate          time.Time
}

type TransactionTypes int

const (
	Transfer TransactionTypes = iota
	Income
	Expense
)
