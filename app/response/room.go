package response

type RooomResponse struct {
	ID              string `bun:"id" json:"id"`
	Name            string `bun:"name" json:"name"`
	RoomTypeID      string `bun:"room_type_id" json:"room_type_id"`
	Capacity        int64  `bun:"capacity" json:"capacity"`
	Description     string `bun:"description" json:"description"`
	ImageURL        string `bun:"image_url" json:"image_url"`
	Building        string `bun:"building" json:"building"`
	StartRoom       int64  `bun:"start_room" json:"start_room"`
	EndRoom         int64  `bun:"end_room" json:"end_room"`
	Is_Available    bool   `bun:"is_available" json:"is_available"`
	MaintenanceNote string `bun:"maintenance_note" json:"maintenance_note"`
	MaintenanceETA  string `bun:"maintenance_eta" json:"maintenance_eta"`
	CreatedAt       int64  `bun:"created_at" json:"created_at"`
	UpdatedAt       int64  `bun:"updated_at" json:"updated_at"`
	DeletedAt       int64  `bun:"deleted_at" json:"deleted_at"`
}
