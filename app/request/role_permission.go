package request

type Role_PermissionCreate struct {
	Role_ID       string `json:"role_id"`
	Permission_ID string `json:"permission_id"`
}

type Role_PermissionUpdate struct {
	Role_PermissionCreate
}

type Role_PermissionGetByID struct {
	ID string `json:"role_id"`
}
