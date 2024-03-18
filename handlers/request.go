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

type GuestReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Date      string `json:"date"`
	Phone     string `json:"phone"`
}

type RequestBody struct {
	Description string     `json:"description"`
	Guests      []GuestReq `json:"guests"`
}

type Request struct {
	db *gorm.DB
}

func NewRequestHandler(db *gorm.DB) *Request {
	return &Request{
		db: db,
	}
}

func (r *Request) CreateRequest(c echo.Context) error {
	reqBody := new(RequestBody)
	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("invalid request body"))
	}

	userID := c.Get("id").(uint)
	guests := []models.Guest{}
	for _, g := range reqBody.Guests {
		date, err := time.Parse("2006/01/02", g.Date)
		if err != nil {
			return InternalServerError(c)
		}

		guest := models.Guest{
			FirstName: g.FirstName,
			LastName:  g.LastName,
			Phone:     g.Phone,
			Date:      date,
		}

		guests = append(guests, guest)
	}

	request := new(models.Request)
	request.Description = reqBody.Description
	request.Guests = guests
	request.UserID = userID
	err := r.db.Create(request).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (r *Request) GetRequests(c echo.Context) error {
	requests := []models.Request{}
	var err error
	cat := c.QueryParam("cat")

	if cat == "sent" {
		err = r.db.Find(&requests, "user_id=? and sent=true", c.Get("id").(uint)).Error
	} else if cat == "unsent" {
		err = r.db.Find(&requests, "user_id=? and sent=false", c.Get("id").(uint)).Error
	} else if cat == "forApproval" {
		err = r.db.Find(&requests, "user_id_for_approval=? and sent=true", c.Get("id").(uint)).Error
	}
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    requests,
	})
}

func (r *Request) GetRequest(c echo.Context) error {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
	}

	request := new(models.Request)
	err = r.db.Preload("Guests").Take(request, uint(reqID)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    request,
	})
}

func (r *Request) DeleteRequest(c echo.Context) error {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Request{}, "id=?", uint(reqID)).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.Guest{}, "request_id=?", uint(reqID)).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (r *Request) SendRequest(c echo.Context) error {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
	}

	userID := c.Get("id").(uint)
	userRole := c.Get("role").(string)

	user := new(models.User)
	err = r.db.Take(user, userID).Error
	if err != nil {
		return InternalServerError(c)
	}

	request := new(models.Request)
	err = r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	if userRole == "user" {
		dean := new(models.User)
		err = r.db.Take(dean, "department_id=? and role=?", user.DepartmentID, "dean").Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return InternalServerError(c)
		}
		if dean.ID == 0 {
			return c.JSON(http.StatusNotFound, util.ErrResp("ریاست این بخش مشخص نشده است"))
		}

		request.UserIdForApproval = &dean.ID
		request.Sent = true
		now := time.Now()
		request.SentDate = &now

		err = r.db.Save(request).Error
		if err != nil {
			return InternalServerError(c)
		}

	} else {
		securityDep := new(models.Department)
		err = r.db.Take(securityDep, "name=?", "حراست کل").Error
		if err != nil {
			return InternalServerError(c)
		}

		if userRole == "securityDean" {
			//request.FinalApproval = true
			approved := true
			request.DepartmentApproval = &approved
			request.SecurityApproval = &approved
			request.Sent = true
			now := time.Now()
			request.SentDate = &now
		} else if userRole == "dean" {
			securityDean := new(models.User)
			err = r.db.Take(securityDean, "role=? and department_id=?", "securityDean", securityDep.ID).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return InternalServerError(c)
			}
			if securityDean.ID == 0 {
				return c.JSON(http.StatusNotFound, util.ErrResp("ریاست حراست مشخص نشده است"))
			}

			approved := true
			request.DepartmentApproval = &approved
			request.Sent = true
			request.UserIdForApproval = &securityDean.ID
			now := time.Now()
			request.SentDate = &now
		}
		err = r.db.Save(request).Error
		if err != nil {
			return InternalServerError(c)
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (r *Request) DeanApproval(c echo.Context) error {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
	}

	userRole := c.Get("role").(string)
	if userRole != "dean" {
		return c.JSON(http.StatusBadRequest, util.ErrResp("این حساب به این امکان دسترسی ندارد"))
	}

	request := new(models.Request)
	err = r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	security := new(models.Department)
	err = r.db.Take(security, "name=?", "حراست کل").Error
	if err != nil {
		return InternalServerError(c)
	}

	securityDean := new(models.User)
	err = r.db.Take(securityDean, "role=? and department_id=?", "securityDean", security.ID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if securityDean.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("ریاست حراست مشخص نشده است"))
	}

	request.UserIdForApproval = &securityDean.ID
	approved := true
	request.DepartmentApproval = &approved
	now := time.Now()
	request.SentDate = &now

	err = r.db.Save(request).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (r *Request) DeanDisapproval(c echo.Context) error {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
	}

	userRole := c.Get("role").(string)
	if userRole != "dean" {
		return c.JSON(http.StatusBadRequest, util.ErrResp("این حساب به این امکان دسترسی ندارد"))
	}

	request := new(models.Request)
	err = r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	disapproved := false
	request.DepartmentApproval = &disapproved
	request.UserIdForApproval = nil

	err = r.db.Save(request).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (r *Request) SecurityApproval(c echo.Context) error {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
	}

	userRole := c.Get("role").(string)
	if userRole != "securityDean" {
		return c.JSON(http.StatusBadRequest, util.ErrResp("این حساب به این امکان دسترسی ندارد"))
	}

	request := new(models.Request)
	err = r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	approved := true
	request.SecurityApproval = &approved
	request.UserIdForApproval = nil
	err = r.db.Save(request).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (r *Request) SecurityDisapproval(c echo.Context) error {
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
	}

	userRole := c.Get("role").(string)
	if userRole != "securityDean" {
		return c.JSON(http.StatusBadRequest, util.ErrResp("این حساب به این امکان دسترسی ندارد"))
	}

	request := new(models.Request)
	err = r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	disapproved := false
	request.SecurityApproval = &disapproved
	request.UserIdForApproval = nil

	err = r.db.Save(request).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}
