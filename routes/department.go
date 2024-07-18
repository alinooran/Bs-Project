package routes

import (
	"github.com/alinooran/Bs-Project/handlers"
	"github.com/alinooran/Bs-Project/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyDepartmentRoutes(g *echo.Group, db *gorm.DB) {
	departmentHandler := handlers.NewDepartmentHandler(db)

	g.GET("", departmentHandler.GetDepartments, middleware.NormalAccess)
}
