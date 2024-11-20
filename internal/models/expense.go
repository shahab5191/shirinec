package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
    ID          int
    UserID      uuid.UUID
    AccountID   int
    CategoryID  int
    Amount      float32
    Description string
    Date        time.Time
}
