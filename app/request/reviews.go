package request

type CreateReviews struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
	Rating int `json:"rating"`
	Comment string `json:"comment"`
}