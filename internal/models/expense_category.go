package models

import (
	"github.com/google/uuid"
)

type ExpenseCategory struct {
	ID     int
	UserID uuid.UUID
	Name   *string
	Color  *string
}
