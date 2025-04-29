package response

type List_Role struct {
	ID   string `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
}
