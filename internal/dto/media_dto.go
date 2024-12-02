package dto

import (
	"time"

	"shirinec.com/internal/enums"
)

type MediaUploadRequest struct {
	BindsTo enums.MediaUploadBind `form:"binds_to,required" binding:"required,mediaUploadBind"`
    BindID  int                   `form:"bind_id,required" binding:"required,number"`
}

type MediaUploadResponse struct {
	ID           int       `json:"id"`
	URL          string    `json:"url"`
	Metadata     *string   `json:"metadata"`
	CreationDate time.Time `json:"creationDate"`
	UpdateDate   time.Time `json:"updateDate"`
}
