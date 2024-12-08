package models

import (
	"time"

	"github.com/google/uuid"
	"shirinec.com/internal/enums"
)

type Media struct {
	ID               int
	UserID           uuid.UUID
	Url              string
	FilePath         string
	Metadata         *string
	Access           *enums.MediaAccess
	FinancialGroupID int
	CreationDate     time.Time
	UpdateDate       time.Time
}
