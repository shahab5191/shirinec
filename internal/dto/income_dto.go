package dto

type ListRequest struct {
	Page int `form:"page,default=0"`
	Size int `form:"size,default=10"`
}
