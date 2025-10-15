package response

type PairingCodeGet struct {
	UserID          string `json:"user_id"`
	Code            string `json:"code"`           
	ExpiresAt       int64  `json:"expires_at"`     
	ExpiresAtHuman  string `json:"expires_at_human"`
	CreatedAt       int64  `json:"created_at"`
}