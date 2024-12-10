package models

import (
	"time"

	"github.com/google/uuid"
)

type FinancialGroups struct {
    ID int
    Name string
    UserID uuid.UUID
    CreationDate time.Time
    UpdateDate time.Time
}
