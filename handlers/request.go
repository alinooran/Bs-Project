package handlers

import (
	"errors"
	"github.com/alinooran/Bs-Project/models"
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Request struct {
	db *gorm.DB
}

func NewRequestHandler(db *gorm.DB) *Request {
	return &Request{
		db: db,
	}
}

func (r *Request) CreateRequest(c echo.Context) error {
	reqBody := &struct {
		Date   string         `json:"date"`
		Guests []models.Guest `json:"guests"`
	}{}

	if err := c.Bind(&reqBody); err != nil {
		return RequestBodyError(c)
	}

	userID := c.Get("id").(uint)
	date, err := time.Parse("2006/01/02", reqBody.Date)
	if err != nil {
		return InternalServerError(c)
	}

	request := new(models.Request)
	request.Guests = reqBody.Guests
	request.Date = date
	request.UserID = userID

	err = r.db.Create(request).Error
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
	userId := c.Get("id").(uint)
	role := c.Get("role").(string)

	if cat == "sent" {
		err = r.db.Order("id desc").Find(&requests, "user_id=? and sent=true", userId).Error
	} else if cat == "unsent" {
		err = r.db.Find(&requests, "user_id=? and sent=false", userId).Error
	} else if cat == "forApproval" {
		err = r.db.Find(&requests, "user_id_for_approval=? and sent=true", userId).Error
	} else {
		var departmentId uint
		err = r.db.Table("users").Select("department_id").Where("id = ?", userId).Scan(&departmentId).Error
		if err != nil {
			return InternalServerError(c)
		}
		if cat == "approved" {
			query := r.db.Table("requests").Select("requests.*").Joins("JOIN users ON requests.user_id = users.id").
				Where("requests.deleted_at IS NULL AND users.department_id = ?", departmentId).Order("updated_at desc")
			if role == "dean" {
				query = query.Where("requests.dean_approval = ?", 1)
			} else if role == "securityDean" {
				query = query.Where("requests.security_approval = ?", 1)
			}
			err = query.Scan(&requests).Error
		} else if cat == "disapproved" {
			query := r.db.Table("requests").Select("requests.*").Joins("JOIN users ON requests.user_id = users.id").
				Where("requests.deleted_at IS NULL AND users.department_id = ?", departmentId).Order("updated_at desc")
			if role == "dean" {
				query = query.Where("requests.dean_approval = ?", 0)
			} else if role == "securityDean" {
				query = query.Where("requests.security_approval = ?", 0)
			}
			err = query.Scan(&requests).Error
		}
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
	reqID := c.Get("request_id").(uint)

	request := new(models.Request)
	err := r.db.Preload("Guests").Take(request, reqID).Error
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
	reqID := c.Get("request_id").(uint)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Request{}, "id=?", reqID).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.Guest{}, "request_id=?", reqID).Error; err != nil {
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
	reqBody := &struct {
		Description string `json:"description"`
	}{}
	if err := c.Bind(reqBody); err != nil {
		return RequestBodyError(c)
	}

	reqID := c.Get("request_id").(uint)

	userID := c.Get("id").(uint)
	userRole := c.Get("role").(string)

	user := new(models.User)
	err := r.db.Take(user, userID).Error
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

	request.Sent = true
	request.FinalStatus = nil

	workflow := &models.Workflow{
		SenderName:  user.FirstName + " " + user.LastName,
		Step:        "ارسال توسط کاربر",
		Description: reqBody.Description,
		RequestID:   reqID,
	}

	approved := true

	if userRole == "user" {
		dean := new(models.User)
		err = r.db.Select("id").Take(dean, "department_id=? and role=?", user.DepartmentID, "dean").Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return InternalServerError(c)
		}
		if dean.ID == 0 {
			return c.JSON(http.StatusNotFound, util.ErrResp("ریاست این بخش مشخص نشده است"))
		}
		if request.UserIdForApproval == nil {
			request.UserIdForApproval = &dean.ID
		}
	} else if userRole == "dean" {
		securityDean := new(models.User)
		err = r.db.Take(securityDean, "role=?", "securityDean").Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return InternalServerError(c)
		}
		if securityDean.ID == 0 {
			return c.JSON(http.StatusNotFound, util.ErrResp("ریاست حراست مشخص نشده است"))
		}
		if request.UserIdForApproval == nil {
			request.UserIdForApproval = &securityDean.ID
		}
		request.DeanApproval = &approved
	} else {
		request.UserIdForApproval = nil
		request.SecurityApproval = &approved
		request.FinalStatus = &approved
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(workflow).Error
		if err != nil {
			return err
		}

		err = tx.Save(request).Error
		if err != nil {
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

func (r *Request) Approve(c echo.Context) error {
	reqID := c.Get("request_id").(uint)

	userRole := c.Get("role").(string)
	if userRole == "user" {
		return c.JSON(http.StatusBadRequest, util.ErrResp("این حساب به این امکان دسترسی ندارد"))
	}

	request := new(models.Request)
	err := r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	userID := c.Get("id").(uint)

	if userID != *request.UserIdForApproval {
		return c.JSON(http.StatusBadRequest, util.ErrResp("کاربر مجاز به انجام این عملیات نیست"))
	}

	user := new(models.User)
	err = r.db.Take(user, userID).Error
	if err != nil {
		return InternalServerError(c)
	}

	approved := true

	workflow := &models.Workflow{
		SenderName: user.FirstName + " " + user.LastName,
		RequestID:  reqID,
		Status:     &approved,
	}

	if userRole == "dean" {
		securityDean := new(models.User)
		err = r.db.Take(securityDean, "role=?", "securityDean").Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return InternalServerError(c)
		}
		if securityDean.ID == 0 {
			return c.JSON(http.StatusNotFound, util.ErrResp("ریاست حراست مشخص نشده است"))
		}

		request.UserIdForApproval = &securityDean.ID
		request.DeanApproval = &approved
		workflow.Step = "تایید ریاست بخش"
	} else {
		request.FinalStatus = &approved
		request.SecurityApproval = &approved
		workflow.Step = "تایید حراست"
		request.UserIdForApproval = nil
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(workflow).Error
		if err != nil {
			return err
		}

		err = tx.Save(request).Error
		if err != nil {
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

func (r *Request) Reject(c echo.Context) error {
	reqBody := &struct {
		Description string `json:"description"`
	}{}
	if err := c.Bind(reqBody); err != nil {
		return RequestBodyError(c)
	}

	reqID := c.Get("request_id").(uint)

	userRole := c.Get("role").(string)
	if userRole == "user" {
		return c.JSON(http.StatusBadRequest, util.ErrResp("این حساب به این امکان دسترسی ندارد"))
	}

	request := new(models.Request)
	err := r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	userID := c.Get("id").(uint)

	if userID != *request.UserIdForApproval {
		return c.JSON(http.StatusBadRequest, util.ErrResp("کاربر مجاز به انجام این عملیات نیست"))
	}

	user := new(models.User)
	err = r.db.Take(user, userID).Error
	if err != nil {
		return InternalServerError(c)
	}

	approved := false

	workflow := &models.Workflow{
		SenderName:  user.FirstName + " " + user.LastName,
		RequestID:   reqID,
		Status:      &approved,
		Description: reqBody.Description,
	}

	if userRole == "dean" {
		request.UserIdForApproval = &userID
		request.DeanApproval = &approved
		workflow.Step = "تایید ریاست بخش"
	} else {
		workflow.Step = "تایید حراست"
		request.UserIdForApproval = &userID
		request.SecurityApproval = &approved
	}
	request.FinalStatus = &approved
	request.Sent = false

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(workflow).Error
		if err != nil {
			return err
		}

		err = tx.Save(request).Error
		if err != nil {
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

func (r *Request) GetWorkflows(c echo.Context) error {
	reqID := c.Get("request_id").(uint)

	workflows := []models.Workflow{}

	err := r.db.Find(&workflows, "request_id=?", reqID).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    workflows,
	})
}

func (r *Request) GetGuests(c echo.Context) error {
	reqID := c.Get("request_id").(uint)

	guests := []models.Guest{}

	err := r.db.Find(&guests, "request_id=?", reqID).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    guests,
	})
}

func (r *Request) CloseRequest(c echo.Context) error {
	reqID := c.Get("request_id").(uint)
	request := new(models.Request)
	err := r.db.Take(request, reqID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}
	if request.ID == 0 {
		return c.JSON(http.StatusNotFound, util.ErrResp("درخواست یافت نشد"))
	}

	request.Sent = true
	err = r.db.Save(request).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, nil)
}
