package response

type Role_PermissionResponse struct {
	ID            string `bun:"id" json:"id"`
	Role_ID       string `bun:"role_id" json:"role_id"`
}