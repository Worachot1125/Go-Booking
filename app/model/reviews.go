package model

import "github.com/uptrace/bun"

type Reviews struct {
	bun.BaseModel `bun:"table:reviews"`
	ID      string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	UserID  *User  `bun:"user_id,notnull"`
	RoomID  *Room  `bun:"room_id,notnull"`
	Rating  int    `bun:"rating,notnull"`
	Comment string `bun:"comment"`

	CreateUpdateUnixTimestamp
	SoftDelete
}