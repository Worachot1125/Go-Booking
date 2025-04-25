package model

import (
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID          string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	FirstName   string `bun:"first_name,notnull"`
	LastName    string `bun:"last_name,notnull"`
	Email       string `bun:"email,unique,notnull"`
	Password    string `bun:"password,notnull"`
	Position_ID string `bun:"position_id"` // Foreign key to the Position table
	Image_url   string `bun:"image_url"`

	Position *Position `bun:"rel:belongs-to,join:position_id=id"`
	CreateUpdateUnixTimestamp
	SoftDelete
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
