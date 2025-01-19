package models

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/src/internal/enums"
)

type Category struct {
	ID           int
	UserID       uuid.UUID
	Name         *string
	Color        *string
	IconID       *int
	EntityType   *enums.CategoryType
	CreationDate *time.Time
	UpdateDate   *time.Time
}
