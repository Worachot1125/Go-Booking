package response

type FindByEmail struct {
	Email string `bun:"email"`
}