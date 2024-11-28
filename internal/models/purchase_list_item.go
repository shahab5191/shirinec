package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseListItem struct {
    ID int
    UserID uuid.UUID
    ItemID int
    Count int
    UnitPrice float64
    TransactionID int
    CreationDate time.Time
    UpdateDate time.Time
}
