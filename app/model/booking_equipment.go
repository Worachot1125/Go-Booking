package model

import "github.com/uptrace/bun"

type BookingEquipment struct {
	bun.BaseModel `bun:"table:building_rooms,alias:br"`

	ID          string     `bun:",pk,type:uuid,default:gen_random_uuid()"`
	BookingID   *Booking   `bun:"booking_id,notnull"`
	EquipmentID *Equipment `bun:"equipment_id,notnull"`
	Quantity    int        `bun:"quantity,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
