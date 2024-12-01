package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID           int
	UserID       uuid.UUID
	Url          string
	FilePath     string
	Metadata     *string
	CreationDate time.Time
	UpdateDate   time.Time
}
