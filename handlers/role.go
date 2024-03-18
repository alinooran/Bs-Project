package handlers

import (
	"github.com/alinooran/Bs-Project/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Role struct {
	db *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *Role {
	return &Role{
		db: db,
	}
}

func (r *Role) GetRoles(c echo.Context) error {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    roles,
	})
}
