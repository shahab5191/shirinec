package models

import "github.com/google/uuid"

type IncomeCategory struct {
	ID     int
	UserID uuid.UUID
	Name   *string
	Color  *string
}
