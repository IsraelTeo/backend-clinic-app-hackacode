package db

import (
	"fmt"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
)

/*
func MigrateDB() error {
	err := GDB.AutoMigrate(
		&model.Service{},
		&model.Package{},
		&model.Patient{},
		&model.User{},
		&model.Appointment{},
		&model.Doctor{},
	)

	if err != nil {
		return err
	}

	return nil
}
*/

// Migrate migrates the entities using GORM
func Migrate(db *Database) error {
	if err := db.DB.AutoMigrate(
		&model.MedicalService{},
		&model.MedicalPackage{},
		&model.Patient{},
		&model.Doctor{},
		&model.User{},
		&model.Appointment{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Database migration completed!")
	return nil
}
