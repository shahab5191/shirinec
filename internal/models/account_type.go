package models

import "github.com/google/uuid"

type AccountType struct {
	ID      int
	user_id uuid.UUID
	name    string
	icon    int
}
