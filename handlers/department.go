package handlers

import (
	"github.com/alinooran/Bs-Project/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Department struct {
	db *gorm.DB
}

func NewDepartmentHandler(db *gorm.DB) *Department {
	return &Department{
		db: db,
	}
}

func (d *Department) GetDepartments(c echo.Context) error {
	var departments []models.Department
	err := d.db.Select("id", "name").Find(&departments).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    departments,
	})
}
