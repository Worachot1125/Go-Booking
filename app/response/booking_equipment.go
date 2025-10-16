package response

type BookingEquipmentResponse struct {
	ID          string `json:"id" bun:"id"`
	BookingID   string `json:"booking_id" bun:"booking_id"`
	EquipmentID string `json:"equipment_id" bun:"equipment_id"`
	Quantity    int    `json:"quantity" bun:"quantity"`
	CreatedAt   int64  `json:"created_at" bun:"created_at"`
	UpdatedAt   int64  `json:"updated_at" bun:"updated_at"`
}
