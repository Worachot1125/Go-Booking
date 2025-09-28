package response

type EquipmentResponse struct {
	ID                 string `json:"id" bun:"id"`
	Name               string `json:"name" bun:"name"`
	Image_URL          string `json:"image_url" bun:"image_url"`
	Quantity           int    `json:"quantity" bun:"quantity"`
	Status             string `json:"status" bun:"status"`
	CreatedAt          int64  `json:"created_at" bun:"created_at"`
	UpdatedAt          int64  `json:"updated_at" bun:"updated_at"`
}

