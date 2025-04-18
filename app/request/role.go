package request

type CreateRole struct {
	Name   string `json:"name"`
}

type GetByIdRole struct {
	ID string `uri:"id" binding:"required"`
}