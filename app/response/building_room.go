package response

type Building_RoomResponse struct{
	ID 		  string `bun:"id"`
	RoomID    string `bun:"room_id"`
	BuildingID string `bun:"building_id"`
	CreatedAt int64  `bun:"created_at"`
	UpdatedAt int64  `bun:"updated_at"`
}