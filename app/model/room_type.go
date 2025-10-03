package model

import "github.com/uptrace/bun"

type RoomType struct {
	bun.BaseModel `bun:"table:room_types,alias:rt"`

	ID   string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name string `bun:"name,unique,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
