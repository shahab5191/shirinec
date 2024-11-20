package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
    ID          uuid.UUID
    Name        string
    Email       string
    IP          string
    Password    string
    LastLogin   time.Time
    FailedTries int
    Status      UserStatus
    Role        UserRole
    CreatedDate time.Time
    UpdatedDate time.Time
}

type UserStatus int

const (
    StatusBanned UserStatus = iota
    StatusValidated
    StatusDisabled
    StatusLocked
    StatusPending
)

type UserRole int

const (
    RoleUser UserRole = iota
    RoleAdmin
)
