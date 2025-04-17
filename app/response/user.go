package response

type ListUser struct{
	ID string `bun:"id"`
	FirstName string `bun:"first_name"`
	LastName string `bun:"last_name"`
	Email string `bun:"email"`
	CreatedAt int64 `bun:"created_at"`
	UpdatedAt int64 `bun:"updated_at"`
}