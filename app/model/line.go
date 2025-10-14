package model

import "github.com/uptrace/bun"

type LinePairingCode struct {
	bun.BaseModel `bun:"table:line_pairing_codes"`
	ID            string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	UserID        string `bun:"user_id,notnull"`
	Code          string `bun:"code,unique,notnull"`
	ExpiresAt     int64  `bun:"expires_at,notnull"`
	UsedAt        int64  `bun:"used_at,nullzero"`

	CreateUpdateUnixTimestamp
}
