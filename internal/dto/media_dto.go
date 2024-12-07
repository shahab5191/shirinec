package dto

import (
	"time"
)

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
