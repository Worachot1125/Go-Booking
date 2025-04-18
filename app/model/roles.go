package model

import (
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID   string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name string `bun:"name,unique,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
