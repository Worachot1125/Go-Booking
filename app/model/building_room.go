package model

import (
	"github.com/uptrace/bun"
)

type Building_Room struct {
	bun.BaseModel `bun:"table:building_rooms,alias:br"`

	ID         string    `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Room 	   *Room 	 `bun:"rel:belongs-to,join:room_id=id"`
	RoomID 	   string    `bun:"room_id,notnull"`
	Building   *Building `bun:"rel:belongs-to,join:building_id=id"`
	BuildingID string 	 `bun:"building_id,notnull"`


	CreateUpdateUnixTimestamp``
	SoftDelete
}
