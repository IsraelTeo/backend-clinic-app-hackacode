package db

import (
	"fmt"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/config"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GDB *gorm.DB

func Connection(cfg *config.Config) error {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	var err error

	GDB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	GDB = GDB.Debug()

	return nil
}

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
