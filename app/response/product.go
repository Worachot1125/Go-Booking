package response

type ProductResponse struct {
	ID          int64   `bun:"id"`
	Name        string  `bun:"name"`
	Price       float64 `bun:"price"`
	Description string  `bun:"description"`
	CreatedAt   string  `bun:"created_at"`
	UpdatedAt   string  `bun:"updated_at"`
}
