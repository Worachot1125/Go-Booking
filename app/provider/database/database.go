package provider

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	dbMap = make(map[string]*bun.DB)
	db    *bun.DB
	db2   *bun.DB
)

type DBOption struct {
	DSN      string
	Host     string
	Port     int64
	Database string
	Username string
	Password string
	TimeZone string
	SSLMode  string
}

// Register registers and initializes the database connection
func Register(conn **bun.DB, conf *DBOption) {
	if conf.DSN == "" {
		conf.DSN = generateDSN(conf)
	}

	config, err := pgx.ParseConfig(conf.DSN)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	sqldb := stdlib.OpenDB(*config)
	*conn = bun.NewDB(sqldb, pgdialect.New())

	if viper.GetBool("DEBUG") {
		(*conn).AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	err = (*conn).Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
}

// generateDSN generates the DSN string for the database connection
func generateDSN(conf *DBOption) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
		conf.Host, conf.Port, conf.Database, conf.Username, conf.Password, conf.SSLMode, conf.TimeZone)
}

// DB returns the main database connection
func DB() *bun.DB {
	return db
}

func DB2() *bun.DB {
	return db2
}

// Open function to open all registered databases
func Open(ctx context.Context) error {
	var err error
	for _, db := range dbMap {
		if errClose := db.Ping(); errClose != nil {
			err = errors.Join(err, errClose)
		}
	}
	return err
}

// Close function to close all registered databases
func Close(ctx context.Context) error {
	var err error
	for _, db := range dbMap {
		if errClose := db.Close(); errClose != nil {
			err = errors.Join(err, errClose)
		}
	}
	return err
}
