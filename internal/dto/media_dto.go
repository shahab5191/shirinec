package dto

import (
	"time"

	"shirinec.com/internal/enums"
)

type MediaUploadRequest struct {
	BindsTo enums.MediaUploadBind `form:"binds_to, required" validate:"required,mediaUploadBind"`
	BindID  int                   `form:"bind_id, required"`
}

type MediaUploadResponse struct {
	ID           int       `json:"id"`
	URL          string    `json:"url"`
	Metadata     *string   `json:"metadata"`
	CreationDate time.Time `json:"creationDate"`
	UpdateDate   time.Time `json:"updateDate"`
}
