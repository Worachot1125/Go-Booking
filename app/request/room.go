package request

type CreateRoom struct {
	Name            string `json:"name" form:"name"`
	RoomTypeID      string `json:"room_type_id" form:"room_type_id"`
	Capacity        int64  `json:"capacity" form:"capacity"`
	Description     string `json:"description" form:"description"`
	Image_url       string `json:"image_url" form:"image_url"`
	StartRoom       int64  `json:"start_room" form:"start_room"`
	EndRoom         int64  `json:"end_room" form:"end_room"`
	Is_Available    bool   `json:"is_available" form:"is_available"`
	MaintenanceNote string `json:"maintenance_note" form:"maintenance_note"`
	MaintenanceETA  string `json:"maintenance_eta" form:"maintenance_eta"`
}

type UpdateRoom struct {
	CreateRoom
	Is_Available *bool `json:"is_available"` // ใช้ pointer เพื่อแยกกรณีไม่ได้ส่ง
}

type ListRoom struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
}

type GetByIdRoom struct {
	ID string `uri:"id" binding:"required"`
}
