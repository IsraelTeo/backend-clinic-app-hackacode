package db

import (
	"fmt"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/config"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GDB *gorm.DB

func Connection(cfg *config.Config) (*gorm.DB, error) {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	GDB, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return GDB, nil
}

func MigrateDB(GDB *gorm.DB) error {
	err := GDB.AutoMigrate(
		&model.Service{},
		&model.Package{},
	)
	if err != nil {
		return err
	}

	return nil
}
