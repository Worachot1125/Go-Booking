package config

import (
	pdb "app/app/provider/database"
	"database/sql"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func Database() {
	dbOnce.Do(func() {
		// ==== ของเดิมทั้งหมดใน Database() ใส่ไว้ในนี้ ====
		if dsn := strings.TrimSpace(os.Getenv("DATABASE_URL")); dsn != "" {
			if u, err := url.Parse(dsn); err == nil {
				log.Printf("[DB] mode=DATABASE_URL host=%s db=%s sslmode=%s",
					u.Hostname(), strings.TrimPrefix(u.Path, "/"), u.Query().Get("sslmode"))
				// เติม sslmode=require อัตโนมัติถ้าขาด
				if u.Query().Get("sslmode") == "" {
					q := u.Query()
					q.Set("sslmode", "require")
					u.RawQuery = q.Encode()
					dsn = u.String()
				}
			} else {
				log.Printf("[DB] mode=DATABASE_URL (warn: parse failed: %v)", err)
			}
			openViaDSN(dsn)
			log.Println("database connected success (via DATABASE_URL)")
			return
		}

		// --- fallback split vars (ของเดิม) ---
		host := confString("DB_HOST", "127.0.0.1")
		port := confInt64("DB_PORT", int64(5432))
		name := confString("DB_DATABASE", "Database")
		user := confString("DB_USER", "postgres")
		ssl := confString("DB_SSLMODE", "disable")
		tz := confString("TZ", "Asia/Bangkok")

		log.Printf("[DB] mode=SPLIT_VARS host=%s port=%d db=%s user=%s sslmode=%s tz=%s",
			host, port, name, user, ssl, tz)

		pdb.Register(&db, &pdb.DBOption{
			Host: host, Port: port, Database: name, Username: user,
			Password: confString("DB_PASSWORD", ""), TimeZone: tz, SSLMode: ssl,
		})

		dbLock.Lock()
		dbMap["default"] = db
		dbLock.Unlock()
		log.Println("database connected success (via split vars)")
	})
}

var (
	db     *bun.DB
	dbMap  = make(map[string]*bun.DB) // Initialize the dbMap
	dbLock sync.RWMutex
	dbOnce sync.Once
)

func openViaDSN(dsn string) {
	sqldb, err := sql.Open("pgx", dsn) // << ใช้ pgx
	if err != nil {
		log.Fatal(err)
	}

	// (optional) ปรับ pool
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(30 * time.Minute)

	if err := sqldb.Ping(); err != nil {
		log.Fatal(err)
	}

	db = bun.NewDB(sqldb, pgdialect.New())
	dbMap["default"] = db
}

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
