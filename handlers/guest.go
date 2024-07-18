package handlers

import (
	"errors"
	"github.com/alinooran/Bs-Project/models"
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

	subQuery := g.db.Select("id").Where("final_status=true and date(date)=?", today).Table("requests")
	err := g.db.Find(&guests, "request_id in (?)", subQuery).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    guests,
	})
}

func (g *Guest) RecordEntry(c echo.Context) error {
	guestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه مهمان معتبر نیست"))
	}

	guest := new(models.Guest)
	err = g.db.Take(guest, uint(guestID)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if guest.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	guest.IsEntered = true
	now := time.Now()
	guest.EnteredDate = &now

	err = g.db.Save(guest).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}
