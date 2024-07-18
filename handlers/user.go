package handlers

import (
	"errors"
	"github.com/alinooran/Bs-Project/models"
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *User {
	return &User{
		db: db,
	}
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
