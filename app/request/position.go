package request

type PositionCreate struct {
	Name string `json:"name"`
}

type PositionUpdate struct {
	PositionCreate
}

type PositionGetByID struct {
	ID string `uri:"id" binding:"required"`
}

type PositionListRequest struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}