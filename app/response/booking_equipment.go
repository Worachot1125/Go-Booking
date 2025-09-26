package response

type BookingEquipmentResponse struct {
	ID string `json:"id" bun:"id"`
	BookingID string `json:"booking_id" bun:"booking_id"`
	EquipmentID string `json:"equipment_id" bun:"equipment_id"`
}