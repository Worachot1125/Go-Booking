package request

type CreateRoom struct {
	Name 		string  `json:"name"`
	Description string  `json:"description"`
	Capacity 	int64 	`json:"capacity"`
	Image_url 	string  `json:"image_url"`
}

type UpdateRoom struct{
	CreateRoom
}

type GetByIdRoom struct{
	ID string `uri:"id" binding:"required"`
}