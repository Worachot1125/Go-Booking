package response

type BookingResponse struct {
	ID           string `json:"id" bun:"id"`
	UserID       string `json:"user_id" bun:"user_id"`
	UserName     string `json:"user_name" bun:"user_name"`
	UserLastName string `json:"user_lastname" bun:"user_lastname"`
	RoomID       string `json:"room_id" bun:"room_id"`
	RoomName     string `json:"room_name" bun:"room_name"`
	Topic        string `json:"topic" bun:"topic"`
	Description  string `json:"description" bun:"description"`
	Capacity     int    `json:"capacity" bun:"capacity"`
	StartTime    string `json:"start_time" bun:"start_time"`
	EndTime      string `json:"end_time" bun:"end_time"`
	Status       string `json:"status" bun:"status"`
	CreatedAt    string `json:"created_at" bun:"created_at"`
	UpdatedAt    string `json:"updated_at" bun:"updated_at"`
}
