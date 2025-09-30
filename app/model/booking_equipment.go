package model

import "github.com/uptrace/bun"

type BookingEquipment struct {
	bun.BaseModel `bun:"table:booking_equipments,alias:be"`

	ID string `bun:",pk,type:uuid,default:gen_random_uuid()"`

	BookingID   string `bun:"booking_id,notnull"`
	EquipmentID string `bun:"equipment_id,notnull"`
	Quantity    int    `bun:"quantity,notnull"`

	Equipment *Equipment `bun:"rel:belongs-to,join:equipment_id=id"`
	Booking   *Booking   `bun:"rel:belongs-to,join:booking_id=id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
