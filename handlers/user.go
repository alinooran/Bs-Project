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
	err = u.db.Create(newUser).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (u *User) GetUsers(c echo.Context) error {
	var users []models.User
	err := u.db.Omit("password").Find(&users).Error
	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
		"data":    users,
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
