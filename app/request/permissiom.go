package request

type PermissionCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionUpdate struct {
	PositionCreate
}

type PermissionGetByID struct {
	ID string `uri:"id" binding:"required"`
}
