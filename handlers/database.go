package handlers

import (
	"net/http"

	"github.com/alinooran/Bs-Project/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

var MODELS = []interface{}{&models.Department{}, &models.Guest{}, &models.Request{}, &models.User{}, &models.Role{}, &models.Workflow{}}

func NewDatabaseHandler(db *gorm.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) CreateDB(c echo.Context) error {
	err := d.db.AutoMigrate(MODELS...)
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, "OK")
}

func (d *Database) DeleteDB(c echo.Context) error {
	_ = d.db.Migrator().DropTable(MODELS...)
	return c.JSON(http.StatusOK, "OK")
}
