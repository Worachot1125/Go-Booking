package request

type CreateRoom struct {
	Name 		string  `json:"name"`
	Description string  `json:"description"`
	Capacity 	int64 	`json:"capacity"`
	Image_url 	string  `json:"image_url"`
}

type UpdateRoom struct{
	CreateRoom
}

type ListRoom struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}

type GetByIdRoom struct{
	ID string `uri:"id" binding:"required"`
}