package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Booking struct {
	bun.BaseModel `bun:"table:bookings"`

	ID          string             `bun:",pk,type:uuid,default:gen_random_uuid()"`
	User        *User              `bun:"rel:belongs-to,join:user_id=id"`
	UserID      string             `bun:"user_id,notnull"`
	Room        *Room              `bun:"rel:belongs-to,join:room_id=id"`
	RoomID      string             `bun:"room_id,notnull"`
	Title       string             `bun:"title,notnull"`
	Description string             `bun:"description,notnull"`
	Phone       string             `bun:"phone,notnull"`
	StartTime   int64              `bun:"start_time,notnull"`
	EndTime     int64              `bun:"end_time,notnull"`
	approvedBy  *User              `bun:"rel:belongs-to,join:user_id=id"`
	ApprovedBy  string             `bun:"approved_by"`
	Status      enum.BookingStatus `bun:"status,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
