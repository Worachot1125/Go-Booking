package response

type ListUser struct {
	ID          string `bun:"id"`
	FirstName   string `bun:"first_name"`
	LastName    string `bun:"last_name"`
	Email       string `bun:"email"`
	Position_ID string `bun:"position_id"`
	Image_url   string `bun:"image_url"`
	Phone       string `bun:"phone"`
	CreatedAt   int64  `bun:"created_at"`
	UpdatedAt   int64  `bun:"updated_at"`
}
