package models

import (
	"time"

	"github.com/google/uuid"
)


type Item struct {
    ID int
    UserID uuid.UUID
    Name *string
    CategoryID *int
    ImageID *int
    CreationDate time.Time
    UpdateDate time.Time
}
