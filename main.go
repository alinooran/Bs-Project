package main

import (
	"github.com/alinooran/Bs-Project/database"
	"github.com/alinooran/Bs-Project/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowCredentials},
		AllowCredentials: true,
	}))
	db := database.GetConn()
	routes.ApplyRoutes(e, db)
	_ = e.Start(":8080")
}
