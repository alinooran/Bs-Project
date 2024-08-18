package handlers

import (
	"errors"
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	db *gorm.DB
}

type ProfileResp struct {
	Id             uint   `json:"id"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	DepartmentId   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
}

func NewUserHandler(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) GetProfile(c echo.Context) error {
	id := c.Get("id").(uint)
	resp := new(ProfileResp)
	query := u.db.Table("users").Select("users.id, CONCAT(users.first_name, ' ', users.last_name) as name, users.role, users.department_id, departments.name as department_name").
		Joins("JOIN departments ON users.department_id = departments.id").Where("users.id = ?", id)
	err := query.Scan(resp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, util.ErrResp("user not found"))
		} else {
			return InternalServerError(c)
		}
	}
	return c.JSON(http.StatusOK, resp)
}
