package model

import (
	"github.com/uptrace/bun"
)

type Room struct {
	bun.BaseModel `bun:"table:rooms"`

	ID          string  `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name  		string  `bun:"name,notnull"`
	Description string  `bun:"description,notnull"`
	Capacity    int64   `bun:"capacity,notnull"`
	Image_url   string  `bun:"image_url,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
