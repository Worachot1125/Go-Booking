package response

type BuildingResponse struct{
	ID 		  string `bun:"id"`
	Name      string `bun:"name"`
	CreatedAt int64  `bun:"created_at"`
	UpdatedAt int64  `bun:"updated_at"`
}