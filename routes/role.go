package routes

import (
	"github.com/alinooran/Bs-Project/handlers"
	"github.com/alinooran/Bs-Project/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyRoleRoutes(g *echo.Group, db *gorm.DB) {
	roleHandler := handlers.NewRoleHandler(db)

	g.GET("", roleHandler.GetRoles, middleware.AdminAccess)
}
