package dto

type ListRequest struct {
	Page int `form:"page,default=0" binding:"number"`
	Size int `form:"size,default=10" binding:"number"`
}
