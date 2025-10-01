package migrations

import "app/app/model"

func Models() []any {
	return []any{
		// (*model.User)(nil),
		//(*model.Booking)(nil),
		//(*model.Building)(nil),
		// (*model.Building_Room)(nil),
		// (*model.Permission)(nil),
		// (*model.Position)(nil),
		// (*model.Role)(nil),
		// (*model.Role_Permission)(nil),
		// (*model.Room)(nil),
		// (*model.User_Role)(nil),
		//(*model.Equipment)(nil),
		//(*model.BookingEquipment)(nil),
		(*model.Reviews)(nil),
		//(*model.Report)(nil),
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
