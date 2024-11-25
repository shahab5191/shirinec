package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
    ID int
    url string
    user_id uuid.UUID
    metadata *string
    CreationDate time.Time
    UpdateDate time.Time
}
