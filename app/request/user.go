package request

type CreateUser struct {
	FirstName     string `json:"first_name" form:"first_name"`
	LastName      string `json:"last_name" form:"last_name"`
	Email         string `json:"email" form:"email"`
	Password      string `json:"password" form:"password"`
	Position_ID   string `json:"position_id" form:"position_id"`
	Position_Name string `json:"position_name" form:"position_name"`
	Image_url     string `json:"image_url" form:"image_url"`
	Phone         string `json:"phone" form:"phone"`
}

type UpdateUser struct {
	CreateUser
}

type ListUser struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIdUser struct {
	ID string `uri:"id" binding:"required"`
}
