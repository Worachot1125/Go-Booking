package request

type ListRoomType struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}

type GetByIdRoomType struct {
	ID string `uri:"id" binding:"required"`
}