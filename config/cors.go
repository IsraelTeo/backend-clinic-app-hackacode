package config

import "github.com/labstack/echo/v4/middleware"

func CorsConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     []string{"https://clinic-administrator.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
	}
}
