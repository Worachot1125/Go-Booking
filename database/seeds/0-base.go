package seeds

import (
	"app/internal/logger"

	"github.com/uptrace/bun"
)

// "app/app/model"

// Seeds Database seeds
func Seeds(db *bun.DB) error {

	seeder := []func(*bun.DB) error{
		mockUpSeed,
		// userSeed,
		// teamSeed,
	}

	for i := range seeder {
		err := seeder[i](db)
		if err != nil {
			logger.Err(i, "\t", err)
		}
	}
	return nil
}
