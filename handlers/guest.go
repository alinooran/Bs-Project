package handlers

import (
	"github.com/alinooran/Bs-Project/models"
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Guest struct {
	db *gorm.DB
}

func NewGuestHandler(db *gorm.DB) *Guest {
	return &Guest{
		db: db,
	}
}

func (g *Guest) GetTodayGuests(c echo.Context) error {
	var guests []models.Guest
	today := time.Now().Format("2006-01-02")

	userRole := c.Get("role").(string)
	if userRole != "security" {
		return c.JSON(http.StatusBadRequest, util.ErrResp("این حساب به این امکان دسترسی ندارد"))
	}

	subQuery := g.db.Select("id").Where("security_approval=true").Table("requests")
	err := g.db.Find(&guests, "date(date)=? and request_id in (?)", today, subQuery).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    guests,
	})
}
