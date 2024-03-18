package routes

import (
	"github.com/alinooran/Bs-Project/handlers"
	"github.com/alinooran/Bs-Project/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyRequestRoutes(g *echo.Group, db *gorm.DB) {
	requestHandler := handlers.NewRequestHandler(db)
	g.Use(middleware.NormalAccess)

	g.POST("", requestHandler.CreateRequest)
	g.GET("", requestHandler.GetRequests)
	g.DELETE("/:id", requestHandler.DeleteRequest)
	g.POST("/:id", requestHandler.SendRequest)
	g.GET("/:id", requestHandler.GetRequest)
	g.POST("/dean_approval/:id", requestHandler.DeanApproval)
	g.POST("/dean_disapproval/:id", requestHandler.DeanDisapproval)
	g.POST("/security_approval/:id", requestHandler.SecurityApproval)
	g.POST("/security_disapproval/:id", requestHandler.SecurityDisapproval)
}
