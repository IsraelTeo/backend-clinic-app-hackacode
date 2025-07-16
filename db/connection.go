package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once     sync.Once
	instance *Database
)

// Database is a struct that holds the GORM DB instance
type Database struct {
	DB *gorm.DB
}

// GetInstance ensures singleton
func GetInstance(cfg *config.Config) *Database {
	once.Do(
		func() {
			db, err := connectDB(cfg)
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}

			instance = &Database{DB: db}
			fmt.Println("¡Database connection established successfully!")
		},
	)

	return instance
}

// connectDB connects to PostgreSQL
func connectDB(cfg *config.Config) (*gorm.DB, error) {
	// Build the connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	// Open a new database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Enable debug mode for detailed logging
	db = db.Debug()

	// Get the underlying *sql.DB from the GORM *gorm.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// Set connection pool settings
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the database is reachable
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

/*
import (
	"fmt"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/config"
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
}*/
