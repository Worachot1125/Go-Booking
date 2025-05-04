package model

import (
	"github.com/uptrace/bun"
)

type User_Role struct {
	bun.BaseModel `bun:"table:user_roles,alias:ur"`

	ID      string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Role_ID string `bun:"role_id,notnull"`
	User_ID string `bun:"user_id,notnull"`

	// ความสัมพันธ์กับ user_role
	Role *Role `bun:"rel:belongs-to,join:role_id=id"`
	User *User `bun:"rel:belongs-to,join:user_id=id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
