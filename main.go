package main

import (
	"fmt"
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/config"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/db"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar las variables de entorno**
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inicializar la configuración cargando las variables de entorno
	cfg := config.InitConfig()

	// Conectar a la base de datos utilizando la configuración cargada
	gdb, err := db.Connection(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("Database connection established successfully!")

	// Migración de entidades
	if err := db.MigrateDB(gdb); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	fmt.Println("Database migration successful")

	// Inicializar servidor Echo
	//e := echo.New()

}
