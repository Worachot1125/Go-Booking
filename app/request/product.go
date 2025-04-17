package request

type ProductCeate struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type ProductUpdate struct {
	ProductCeate
}

type ProductGetByID struct {
	ID int64 `uri:"id" binding:"required"`
}

type ProductListReuest struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}
