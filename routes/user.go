package routes

import (
	"github.com/alinooran/Bs-Project/handlers"
	"github.com/alinooran/Bs-Project/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyUserRoutes(g *echo.Group, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)

	g.POST("", userHandler.CreateUser, middleware.AdminAccess)
	g.GET("", userHandler.GetUsers, middleware.AdminAccess)
	g.GET("/profile", userHandler.GetProfile, middleware.NormalAccess)
}
