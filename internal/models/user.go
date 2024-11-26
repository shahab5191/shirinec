package models

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/enums"
)

type User struct {
	ID                 uuid.UUID
	Email              string
	IP                 string
	Password           string
	LastLogin          time.Time
	LastPasswordChange time.Time
	FailedTries        int
	Status             enums.UserStatus
	CreationDate       time.Time
	UpdateDate         time.Time
	ProfileID          int
}
