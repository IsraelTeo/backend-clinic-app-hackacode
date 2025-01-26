package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Se define una estructura llamada Config que agrupa varios de configuración
type Config struct {
	PublicHost            string //Almacena la dirección pública del servidor
	Port                  string //Almacena el puerto en el que se ejecutará la aplicación
	DBUser                string //Almacena el nombre de usuario de la base de datos
	DBPassword            string //Almacena la contraseña de la base de datos.
	DBHost                string //Almacena la dirección del servidor de la base de datos
	DBPort                string //Almecena el puerto de db
	DBName                string //Almacena el nombre de la base de datos.
	JWTExpirationInSecond int64  //Almacena la duración de expiración del token JWT en segundos.
	JWTSecret             string //Almacena la clave secreta para los tokens JWT.
}

var Envs = InitConfig()

func InitConfig() *Config {
	jwtExp, err := strconv.ParseInt(os.Getenv("JWT_EXP"), 10, 64)
	if err != nil {
		log.Printf("Error converting JWT_EXP: %v. The default value of 1 hour (%d seconds) will be used", err, jwtExp)
		jwtExp = 3600
	}

	return &Config{
		PublicHost:            os.Getenv("PUBLIC_HOST"),
		Port:                  os.Getenv("PORT"),
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		DBName:                os.Getenv("DB_NAME"),
		JWTExpirationInSecond: jwtExp,
		JWTSecret:             os.Getenv("API_SECRET"),
	}
}

func StartServer(e *echo.Echo, port string) error {
	fmt.Printf("Server starting on port: %s...\n", port)

	if err := e.Start(port); err != nil {
		return fmt.Errorf("error starting server on port %s: %w", port, err)
	}

	return nil
}
