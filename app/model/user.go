package model

import (
	"database/sql"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           string         `bun:",pk,type:uuid,default:gen_random_uuid()"`
	FirstName    string         `bun:"first_name,notnull"`
	LastName     string         `bun:"last_name,notnull"`
	Email        string         `bun:"email,unique,notnull"`
	Password     string         `bun:"password,notnull"`
	Position_ID  string         `bun:"position_id"`
	Image_url    string         `bun:"image_url"`
	Phone        string         `bun:"phone"`
	LineUserID   sql.NullString `bun:"line_user_id,nullzero"`
	LineOptIn    bool           `bun:"line_opt_in"`
	LineLinkedAt sql.NullInt64  `bun:"line_linked_at,nullzero"`

	Position *Position `bun:"rel:belongs-to,join:position_id=id"`
	CreateUpdateUnixTimestamp
	SoftDelete
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
