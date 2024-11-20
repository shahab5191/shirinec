package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
    ID          int
    UserID      uuid.UUID
    Name        string
    Type        AccountTypes
    Balance     float32
    CreatedAt   time.Time
}

type AccountTypes int

const (
    Cash AccountTypes = iota
    CreditCard
    Bank
)
