package request

type CreateReport struct {
	UserID      string `json:"user_id"`
	RoomID      string `json:"room_id"`
	Description string `json:"description"`
}

type UpdateReport struct {
	CreateReport
}

type ListReport struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDReport struct {
	ID string `uri:"id" binding:"required"`
}