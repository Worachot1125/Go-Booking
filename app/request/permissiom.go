package request

type PermissionCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionUpdate struct {
	PermissionCreate
}

type PermissionGetByID struct {
	ID string `uri:"id" binding:"required"`
}
