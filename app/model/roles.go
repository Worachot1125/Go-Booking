package model

import (
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:users"`

	ID   string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name string `bun:"name,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
