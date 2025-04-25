package request

type User_RoleCreate struct {
	User_ID string `json:"user_id" binding:"required"`
	Role_ID string `json:"role_id" binding:"required"`
}