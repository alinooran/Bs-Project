package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ApplyRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api")

	ApplyDBRoutes(api.Group("/db"), db)
	ApplyAuthRoutes(api.Group("/auth"), db)
	ApplyUserRoutes(api.Group("/user"), db)
	ApplyRoleRoutes(api.Group("/role"), db)
	ApplyDepartmentRoutes(api.Group("/department"), db)
	ApplyRequestRoutes(api.Group("/request"), db)
	ApplyGuestRoutes(api.Group("/guest"), db)
}
