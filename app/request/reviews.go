package request

type CreateReviews struct {
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id"`
	BookingID string `json:"booking_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

type UpdateReviews struct {
	CreateReviews
}

type ListReviews struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDReviews struct {
	ID string `uri:"id" binding:"required"`
}

type GetByBookingIDReviews struct {
	BookingID string `uri:"booking_id" binding:"required"`
}