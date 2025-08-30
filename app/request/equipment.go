package request

type CreateEquipment struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Image_URL          string `json:"image_url"`
	Quantity           int    `json:"quantity"`
	Available_Quantity int    `json:"available_quantity"`
	Status             string `json:"status"`
}

type UpdateEquipment struct {
	CreateEquipment
}

type ListEquipment struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	Status   string `form:"status"`
}

type GetByIdEquipment struct {
	ID string `uri:"id" binding:"required"`
}
