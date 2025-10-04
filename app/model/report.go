package model

import "github.com/uptrace/bun"

type Report struct {
	bun.BaseModel `bun:"table:reports,alias:rep"`

	ID          string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	UserID      string `bun:"user_id,notnull"`
	Name_user   string `bun:"name_user,scanonly"`
	
	RoomID      string `bun:"room_id,notnull"`
	Description string `bun:"description"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
	Room *Room `bun:"rel:belongs-to,join:room_id=id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
