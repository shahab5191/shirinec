package dto

type PaginationData struct {
	PageNumber     int `json:"pageNumber"`
	PageSize       int `json:"pageSize"`
	TotalRecord    int `json:"totalRecord"`
	RemainingPages int `json:"remainingPages"`
}
