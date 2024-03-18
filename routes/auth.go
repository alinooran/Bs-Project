package routes

import (
	"github.com/alinooran/Bs-Project/handlers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyAuthRoutes(g *echo.Group, db *gorm.DB) {
	authHandler := handlers.NewAuthHandler(db)

	g.POST("/login", authHandler.Login)
	g.POST("/logout", authHandler.Logout)
}
