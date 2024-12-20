package dto

type PaginationData struct {
	PageNumber     int `json:"pageNumber"`
	PageSize       int `json:"pageSize"`
	TotalRecord    int `json:"totalRecord"`
	RemainingPages int `json:"remainingPages"`
}

type CreateResponse[T any] struct {
	Result T `json:"result"`
}

type ListRequest struct {
	Page int `form:"page,default=0" binding:"number"`
	Size int `form:"size,default=10" binding:"number"`
}
