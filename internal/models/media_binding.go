package models

import (
	"shirinec.com/internal/enums"
)

type MediaAssosiation struct {
	ID          int                   `json:"id"`
	MediaID     int                   `json:"mediaID"`
	BindingType enums.MediaBindType `json:"bindingType"`
	BindingID   int                   `json:"bindingID"`
}
