package model

import (
	"github.com/uptrace/bun"
)

type Room struct {
	bun.BaseModel `bun:"table:rooms"`

	ID               string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	RoomTypeID       string `bun:"room_type_id,type:uuid,notnull"`
	Name             string `bun:"name,notnull"`
	Capacity         int64  `bun:"capacity,notnull"`
	Description      string `bun:"description,notnull"`
	Image_url        string `bun:"image_url,notnull"`
	StartRoom        int64  `bun:"start_room,notnull"`
	EndRoom          int64  `bun:"end_room,notnull"`
	Is_Available     bool   `bun:"is_available,default:true"`
	Maintenance_note string `bun:"maintenance_note"`
	Maintenance_eta  string `bun:"maintenance_eta"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
