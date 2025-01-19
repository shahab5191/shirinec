package dto

import (
	"time"

	"shirinec.com/src/internal/enums"
)

type MediaUploadQuery struct {
	Access           enums.MediaAccess `form:"access" binding:"omitempty,alpha"`
	FinancialGroupID int               `form:"group" binding:"required,number"`
}

type MediaUploadResponse struct {
	ID           int       `json:"id"`
	URL          string    `json:"url"`
	Metadata     *string   `json:"metadata"`
	CreationDate time.Time `json:"creationDate"`
	UpdateDate   time.Time `json:"updateDate"`
}

type MediaListForCleanupResult struct {
	ID       int
	FilePath string
}
