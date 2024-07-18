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

	request := g.Group("/:id", middleware.ParseRequestID)

	request.DELETE("", requestHandler.DeleteRequest)
	request.POST("", requestHandler.SendRequest)
	request.GET("", requestHandler.GetRequest)
	request.GET("/workflow", requestHandler.GetWorkflows)
	request.GET("/guest", requestHandler.GetGuests)
	request.POST("/approve", requestHandler.Approve)
	request.POST("/reject", requestHandler.Reject)
	request.POST("/close", requestHandler.CloseRequest)
}
