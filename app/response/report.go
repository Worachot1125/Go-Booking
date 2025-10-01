package response

type ReportResponse struct {
	ID          string `bun:"id" json:"id"`
	User_ID     string `bun:"user_id" json:"user_id"`
	Room_ID     string `bun:"room_id" json:"room_id"`
	Description string `bun:"description" json:"description"`
	CreatedAt   int64  `bun:"created_at" json:"created_at"`
	UpdatedAt   int64  `bun:"updated_at" json:"updated_at"`
}
