package response

import "encoding/json"

type BuildingResponse struct {
	ID        string          `bun:"id" json:"id"`
	Name      string          `bun:"name" json:"name"`
	RoomsName json.RawMessage `json:"rooms_name"`
	CreatedAt int64           `bun:"created_at" json:"created_at"`
	UpdatedAt int64           `bun:"updated_at" json:"updated_at"`
}
