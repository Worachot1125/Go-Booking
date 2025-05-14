package response

type ListUser struct {
	ID           string `bun:"id" json:"id"`
	FirstName    string `bun:"first_name" json:"first_name"`
	LastName     string `bun:"last_name" json:"last_name"`
	Email        string `bun:"email" json:"email"`
	Position_ID  string `bun:"position_id" json:"position_id"`
	PositionName string `bun:"position_name" json:"position_name"`
	Image_url    string `bun:"image_url" json:"image_url"`
	Phone        string `bun:"phone" json:"phone"`
	CreatedAt    int64  `bun:"created_at" json:"created_at"`
	UpdatedAt    int64  `bun:"updated_at" json:"updated_at"`
}
