package model

import "github.com/uptrace/bun"

type Reviews struct {
	bun.BaseModel `bun:"table:reviews"`
	ID            string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	UserID        string `bun:"user_id,notnull"`
	RoomID        string `bun:"room_id,notnull"`
	Rating        int    `bun:"rating,notnull"`
	Comment       string `bun:"comment"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
	Room *Room `bun:"rel:belongs-to,join:room_id=id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
