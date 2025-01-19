package models

import (
	"time"

	"github.com/google/uuid"
)

type FinancialGroups struct {
    ID int
    Name string
    UserID uuid.UUID
    ImageID *int
    CreationDate time.Time
    UpdateDate time.Time
}
