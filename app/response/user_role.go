package response

type User_RoleResponse struct {
	User_ID string `bun:"user_id"`
	Role_ID string `bun:"role_id"`
}
