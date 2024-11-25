package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	IP           string
	Password     string
	LastLogin    time.Time
	FailedTries  int
	Status       UserStatus
	CreationDate time.Time
	UpdateDate   time.Time
	ProfileID      int
}

type UserStatus int

const (
	StatusBanned UserStatus = iota
	StatusValidated
	StatusDisabled
	StatusLocked
	StatusPending
)
