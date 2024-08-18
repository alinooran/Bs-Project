package routes

import (
	"github.com/alinooran/Bs-Project/handlers"
	"github.com/alinooran/Bs-Project/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyGuestRoutes(g *echo.Group, db *gorm.DB) {
	guestHandler := handlers.NewGuestHandler(db)
	g.Use(middleware.NormalAccess)

	g.GET("", guestHandler.GetTodayGuests)
	g.POST("/:id", guestHandler.RecordEntry)
	g.POST("/report", guestHandler.GetReport)
}
