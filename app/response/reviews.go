package response

type ReviewsResponse struct {
	ID         string `bun:"id" json:"id"`
	User_ID    string `bun:"user_id" json:"user_id"`
	Room_ID    string `bun:"room_id" json:"room_id"`
	Booking_ID string `bun:"booking_id" json:"booking_id"`
	Rating     int64  `bun:"rating" json:"rating"`
	Comment    string `bun:"comment" json:"comment"`
	CreatedAt  int64  `bun:"created_at" json:"created_at"`
	UpdatedAt  int64  `bun:"updated_at" json:"updated_at"`
}
