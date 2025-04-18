package model

import (
	"github.com/uptrace/bun"
)

type Permission struct {
	bun.BaseModel `bun:"table:permissions,alias:p"`

	ID          string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name        string `bun:"name,unique,notnull"`
	Description string `bun:"description"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
