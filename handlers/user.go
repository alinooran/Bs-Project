package handlers

import (
	"errors"
	"github.com/alinooran/Bs-Project/models"
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type SignupReq struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	DepartmentID string `json:"department"`
}

type User struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) CreateUser(c echo.Context) error {
	reqBody := new(SignupReq)
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("invalid request body"))
	}

	dbUser := new(models.User)
	err := u.db.Take(dbUser, "username=?", reqBody.Username).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}

	if dbUser.ID != 0 {
		return c.JSON(http.StatusBadRequest, util.ErrResp("نام کاربری تکراری است"))
	}

	//dbUser := new(models.User)
	err = u.db.Take(dbUser, "email=?", reqBody.Email).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}

	if dbUser.ID != 0 {
		return c.JSON(http.StatusBadRequest, util.ErrResp("ایمیل وارد شده از قبل در سامانه ثبت شده است"))
	}

	departmentId, err := strconv.Atoi(reqBody.DepartmentID)
	if err != nil {
		return InternalServerError(c)
	}

	hashedPass, err := util.HashPassword(reqBody.Password)
	if err != nil {
		return InternalServerError(c)
	}

	securityDep := new(models.Department)
	err = u.db.Take(securityDep, "name=?", "حراست کل").Error
	if err != nil {
		return InternalServerError(c)
	}

	role := reqBody.Role
	if securityDep.ID == uint(departmentId) && role == "dean" {
		role = "securityDean"
	}

	newUser := &models.User{
		Username:     reqBody.Username,
		Password:     hashedPass,
		Role:         role,
		FirstName:    reqBody.FirstName,
		LastName:     reqBody.LastName,
		Phone:        reqBody.Phone,
		Email:        reqBody.Email,
		DepartmentID: uint(departmentId),
	}

	err = u.db.Transaction(func(tx *gorm.DB) error {
		if err := u.db.Create(newUser).Error; err != nil {
			return err
		}
		if err := util.SendEmail(reqBody.Email, reqBody.Username, reqBody.Password); err != nil {
			return err
		}
		return nil
	})

	//err = u.db.Create(newUser).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

//func (u *User) GetUsers(c echo.Context) error {
//	var users []models.User
//	err := u.db.Omit("password").Find(&users).Error
//	if err != nil {
//		return InternalServerError(c)
//	}
//
//	return c.JSON(http.StatusOK, echo.Map{
//		"message": "OK",
//		"data":    users,
//	})

func (u *User) GetUser(c echo.Context) error {
	userId := c.Get("id").(uint)
	dbUser := new(models.User)
	err := u.db.Select("username", "first_name", "last_name", "phone", "email").Take(dbUser, userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, util.ErrResp("کاربر مورد نظر یافت نشد"))
		} else {
			return InternalServerError(c)
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data": echo.Map{
			"username":   dbUser.Username,
			"first_name": dbUser.FirstName,
			"last_name":  dbUser.LastName,
			"phone":      dbUser.Phone,
			"email":      dbUser.Email,
		},
	})
}

func (u *User) GetProfile(c echo.Context) error {
	id := c.Get("id").(uint)
	var dbUser models.User
	err := u.db.Select("first_name", "last_name", "role").Take(&dbUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, util.ErrResp("user not found"))
		} else {
			return InternalServerError(c)
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"name": dbUser.FirstName + " " + dbUser.LastName,
		"role": dbUser.Role,
	})
}

func (u *User) EditPassword(c echo.Context) error {
	reqBody := &struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("invalid request body"))
	}

	userId := c.Get("id").(uint)
	dbUser := new(models.User)
	err := u.db.Take(dbUser, userId).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return InternalServerError(c)
	}

	if dbUser.ID == 0 {
		return c.JSON(http.StatusBadRequest, util.ErrResp("کاربر مورد نظر یافت نشد"))
	}

	if !util.VerifyPassword(dbUser.Password, reqBody.OldPassword) {
		return c.JSON(http.StatusBadRequest, util.ErrResp("رمز عبور وارد شده صحیح نیست"))
	}

	newPassword, err := util.HashPassword(reqBody.NewPassword)
	if err != nil {
		return InternalServerError(c)
	}

	err = u.db.Model(dbUser).Update("password", newPassword).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (u *User) EditUser(c echo.Context) error {
	userId := c.Get("id").(uint)
	reqBody := new(models.User)
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("invalid request body"))
	}

	dbUser := new(models.User)
	err := u.db.Take(dbUser, userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, util.ErrResp("کاربر مورد نظر یافت نشد"))
		} else {
			return InternalServerError(c)
		}
	}

	err = u.db.Take(&models.User{}, "id <> ? AND username = ?", userId, reqBody.Username).Error
	exists := true
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			exists = false
		} else {
			return InternalServerError(c)
		}
	}
	if exists {
		return c.JSON(http.StatusBadRequest, util.ErrResp("نام کاربری تکراری است"))
	}

	err = u.db.Take(&models.User{}, "id <> ? AND email = ?", userId, reqBody.Email).Error
	exists = true
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			exists = false
		} else {
			return InternalServerError(c)
		}
	}
	if exists {
		return c.JSON(http.StatusBadRequest, util.ErrResp("ایمیل وارد شده از قبل در سامانه ثبت شده است"))
	}

	err = u.db.Model(dbUser).Updates(reqBody).Error

	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}
