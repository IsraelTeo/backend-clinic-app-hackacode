package main

import (
	"fmt"
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/config"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/db"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/routes"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/validate"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Cargar las variables de entorno
	godotenv.Load()

	// Inicializar la configuración cargando las variables de entorno
	cfg := config.InitConfig()

	// Paso 1: obtener instancia (conecta y guarda singleton)
	dbInstance := db.GetInstance(cfg)

	// Paso 2: migrar
	if err := db.Migrate(dbInstance); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	fmt.Println("🚀 App started successfully")

	// Inicializar servidor Echo
	e := echo.New()

	//Asignar el validador a la instancia de Echo
	e.Validator = validate.Init()

	//Instanciar Rutas
	routes.InitEnpoints(e)

	//Middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Habilitar CORS
	e.Use(middleware.CORSWithConfig(config.CorsConfig()))

	//Inicia servidor en el puerto: 8080
	err = config.StartServer(e, ":8080")
	if err != nil {
		log.Fatalf("err: %v", err)
	}

}
