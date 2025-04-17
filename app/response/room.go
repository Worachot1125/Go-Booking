package response

type RooomResponse struct {
	ID          string `bun:"id"`
	Name        string `bun:"name"`
	Description string `bun:"description"`
	Capacity    int64  `bun:"capacity"`
	Image_url   string `bun:"image_url"`
	CreatedAt   string `bun:"created_at"`
	UpdatedAt   string `bun:"updated_at"`
}