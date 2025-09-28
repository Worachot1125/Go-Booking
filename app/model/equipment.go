package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Equipment struct {
	bun.BaseModel `bun:"table:equipments,alias:br"`

	ID                 string               `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name               string               `bun:"name,notnull"`
	Image_URL          string               `bun:"image_url,notnull"`
	Quantity           int                  `bun:"quantity,notnull"`
	Status             enum.EquipmentStatus `bun:"status"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
