package migrations

import "app/app/model"

func Models() []any {
	return []any{
		(*model.Building_Room)(nil),
	}
}

func RawBeforeQueryMigrate() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	}
}

func RawAfterQueryMigrate() []string {
	return []string{}
}
