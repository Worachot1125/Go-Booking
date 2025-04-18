package model

import (
	"github.com/uptrace/bun"
)

type Position struct {
	bun.BaseModel `bun:"table:positions"`

	ID   string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name string `bun:"name,unique,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
