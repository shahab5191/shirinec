package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
    ID int
    user_id uuid.UUID
    name string
    AccountTypeID int
    Balance float64
    CreationDate time.Time
    UpdateDate time.Time
}
