package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           int
	UserID       uuid.UUID
	Name         *string
	CategoryID   *int
	Balance      *float64
	CreationDate time.Time
	UpdateDate   time.Time
}
