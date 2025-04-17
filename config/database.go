package config

import (
	pdb "app/app/provider/database"
	"log"
	"sync"

	"github.com/uptrace/bun"
)

func Database() {
	// Connect to database
	pdb.Register(
		&db,
		&pdb.DBOption{
			Host:     confString("DB_HOST", "127.0.0.1"),
			Port:     confInt64("DB_PORT", int64(5432)),
			Database: confString("DB_DATABASE", "Database"),
			Username: confString("DB_USER", "postgres"),
			Password: confString("DB_PASSWORD", ""),
			TimeZone: confString("TZ", "Asia/Bangkok"),
			SSLMode:  confString("DB_SSLMODE", "disable"),
		},
	)
	log.Println("database connected success")

}

var (
	db     *bun.DB
	dbMap  = make(map[string]*bun.DB) // Initialize the dbMap
	dbLock sync.RWMutex
)

func GetDB() *bun.DB {
	return db
}

func DB(name ...string) *bun.DB {
	dbLock.RLock()
	defer dbLock.RUnlock()
	if dbMap == nil {
		panic("database not initialized") // Panic if dbMap is nil
	}
	if len(name) == 0 {
		return dbMap["default"] // Return the default database
	}

	db, ok := dbMap[name[0]]
	if !ok {
		panic("database not initialized") // Panic if the specified database is not found
	}
	return db
}

// func GetDB2() *bun.DB {
// 	return db2
// }
