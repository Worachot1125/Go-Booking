package database

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Option struct {
	DSN      string
	Host     string
	Port     int64
	Database string
	Username string
	Password string
	TimeZone string
}

func New(opt *Option) (*bun.DB, error) {
	if opt.DSN == "" {
		opt.DSN = generateDSN(opt)
	}
	config, err := pgx.ParseConfig(opt.DSN)
	if err != nil {
		return nil, err
	}
	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	return db, db.Ping()
}

func generateDSN(opt *Option) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=%s",
		opt.Host, opt.Port, opt.Database, opt.Username, opt.Password, opt.TimeZone)
}
