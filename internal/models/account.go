package models

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/enums"
)

type Account struct {
	ID           int
	UserID       uuid.UUID
	Name         *string
	CategoryID   *int
	Balance      *float64
	Type         *enums.AccountType
	CreationDate time.Time
	UpdateDate   time.Time
}
