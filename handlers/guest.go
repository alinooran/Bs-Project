package handlers

import (
	"errors"
	"github.com/alinooran/Bs-Project/models"
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Guest struct {
	db *gorm.DB
}

type ReportResp struct {
	ID             uint       `json:"id"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	NationalCode   string     `json:"national_code"`
	EntranceTime   *time.Time `json:"entrance_time"`
	DepartmentName string     `json:"department_name"`
	HostName       string     `json:"host_name"`
	IsEntered      bool       `json:"is_entered"`
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
		return c.JSON(http.StatusNotFound, util.ErrResp("مهمان یافت نشد"))
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

func (g *Guest) GetReport(c echo.Context) error {
	reqBody := &struct {
		Date         string `json:"date"`
		DepartmentID string `json:"department_id"`
		GuestName    string `json:"guest_name"`
		IsEntered    string `json:"is_entered"`
		HostName     string `json:"host_name"`
	}{}
	if err := c.Bind(&reqBody); err != nil {
		log.Printf("[ERR] %#v", err.Error())
		return RequestBodyError(c)
	}

	departmentId, _ := strconv.Atoi(reqBody.DepartmentID)
	departmentIdUint := uint(departmentId)
	reportResp := []ReportResp{}

	query := g.db.Table("guests").Select("guests.id, guests.first_name, guests.last_name, guests.national_code, guests.entered_date as entrance_time, departments.name as department_name, concat(users.first_name, ' ', users.last_name) as host_name, guests.is_entered").
		Joins("JOIN requests ON requests.id = guests.request_id").
		Joins("JOIN users ON users.id = requests.user_id").
		Joins("JOIN departments ON users.department_id = departments.id")

	date, err := time.Parse("2006/01/02", reqBody.Date)
	if err != nil {
		//log.Println("[ERR] date parse error")
		//log.Println(err.Error())
		return InternalServerError(c)
	}

	query = query.Where("requests.date = ?", date).Where("requests.final_status = ?", 1)

	if departmentIdUint != 0 {
		query = query.Where("departments.id = ?", departmentIdUint)
	}

	if reqBody.GuestName != "" {
		query = query.Where("concat(guests.first_name, ' ', guests.last_name) = ?", reqBody.GuestName)
	}

	if reqBody.IsEntered != "" {
		entered := 0
		if reqBody.IsEntered == "true" {
			entered = 1
		}
		query = query.Where("guests.is_entered = ?", entered)
	}

	if reqBody.HostName != "" {
		query = query.Where("concat(users.first_name, ' ', users.last_name) = ?", reqBody.HostName)
	}

	err = query.Scan(&reportResp).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    reportResp,
	})
}
