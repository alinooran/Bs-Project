package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/alinooran/Bs-Project/models"
	"github.com/alinooran/Bs-Project/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *Auth {
	return &Auth{
		db: db,
	}
}

func (a *Auth) Login(c echo.Context) error {
	reqBody := new(LoginReq)
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, util.ErrResp("invalid request body"))
	}

	dbUser := new(models.User)
	err := a.db.Where("username=?", reqBody.Username).Take(dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, util.ErrResp("نام کاربری یا رمز عبور اشتباه است"))
		} else {
			return InternalServerError(c)
		}
	}

	if ok := util.VerifyPassword(dbUser.Password, reqBody.Password); !ok {
		return c.JSON(http.StatusBadRequest, util.ErrResp("نام کاربری یا رمز عبور اشتباه است"))
	}

	claims := jwt.MapClaims{
		"id":   dbUser.ID,
		"role": dbUser.Role,
	}
	token, err := util.GenerateToken(claims)
	if err != nil {
		return InternalServerError(c)
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(2 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"role": dbUser.Role,
	})
}

func (a *Auth) Logout(c echo.Context) error {
	_, err := c.Cookie("token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, util.ErrResp("you must login first"))
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "logout was successful",
	})
}
