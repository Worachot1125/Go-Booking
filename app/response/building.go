package response

type BuildingResponse struct {
	ID        string         `bun:"id" json:"id"`
	Name      string         `bun:"name" json:"name"`
	Image_url string         `bun:"image_url" json:"image_url"`
	Rooms     []RooomResponse `json:"rooms"` // เปลี่ยนตรงนี้
	CreatedAt int64          `bun:"created_at" json:"created_at"`
	UpdatedAt int64          `bun:"updated_at" json:"updated_at"`
}
