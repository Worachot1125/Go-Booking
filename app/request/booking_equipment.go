package request

type CreateBookingEquipment struct {
	BookingID   string `json:"booking_id"`
	EquipmentID string `json:"equipment_id"`
}

type UpdateBookingEquipment struct {
	CreateBookingEquipment
}

type ListEquipmentBooking struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}

type EquipmentSelection struct {
    EquipmentID string `json:"equipment_id" form:"equipment_id"`
    Quantity    int    `json:"quantity" form:"quantity"`
}