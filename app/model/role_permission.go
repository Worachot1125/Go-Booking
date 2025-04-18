package model

import (
	"github.com/uptrace/bun"
)

type Role_Permission struct {
	bun.BaseModel `bun:"table:role_permissions,alias:rp"`

	ID            string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Role_ID       string `bun:"role_id,notnull"`
	Permission_ID string `bun:"permission_id,notnull"`

	// ความสัมพันธ์กับ Role_Permission
	Role       *Role       `bun:"rel:belongs-to,join:role_id=id"`
	Permission *Permission `bun:"rel:belongs-to,join:permission_id=id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
