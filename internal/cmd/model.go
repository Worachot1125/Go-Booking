package cmd

import (
	"context"
	"app/database/migrations"
	"app/database/seeds"
	"app/internal/logger"

	"github.com/uptrace/bun"
)

func modelUp(db *bun.DB) error {
	logger.Infof("Executing model up...")
	for _, mod := range migrations.Models() {
		if _, err := db.NewCreateTable().Model(mod).Exec(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func modelDown(db *bun.DB) error {
	logger.Infof("Executing model down...")
	for _, mod := range migrations.Models() {
		if _, err := db.NewDropTable().Model(mod).Exec(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func modelSeed(db *bun.DB) error {
	logger.Infof("Executing model seeding...")
	if err := seeds.Seeds(db); err != nil {
		logger.Err(err)
	}
	return nil
}

// func modelRawBeforeQuery(db *bun.DB) error {
// 	logger.Infof("Executing pre raw query...")
// 	for _, mod := range migrations.RawBeforeQueryMigrate() {
// 		_, err := db.NewRaw(mod).Exec(context.Background())
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func modelRawAfterQuery(db *bun.DB) error {
// 	logger.Infof("Executing post raw query...")
// 	for _, mod := range migrations.RawAfterQueryMigrate() {
// 		_, err := db.NewRaw(mod).Exec(context.Background())
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
