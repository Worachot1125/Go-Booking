package response

type List_User_Role struct {
	User_ID string `bun:"user_id" json:"user_id"`
	Role_ID string `bun:"role_id" json:"role_id"`
}
