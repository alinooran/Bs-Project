package routes

import (
	"github.com/alinooran/Bs-Project/handlers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyDBRoutes(g *echo.Group, db *gorm.DB) {
	dbHandler := handlers.NewDatabaseHandler(db)

	g.GET("/create", dbHandler.CreateDB)
	g.GET("/delete", dbHandler.DeleteDB)
}
