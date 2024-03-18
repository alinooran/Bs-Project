package middleware

import (
	"github.com/alinooran/Bs-Project/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AdminAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, util.ErrResp("you are not logged in"))
		}

		claims, err := util.ParseToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, util.ErrResp("you dont have access"))
		}

		claimsMap := claims.(jwt.MapClaims)
		if claimsMap["role"] != "admin" {
			return c.JSON(http.StatusUnauthorized, util.ErrResp("you dont have access"))
		}
		return next(c)
	}
}

func NormalAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, util.ErrResp("you are not logged in"))
		}

		claims, err := util.ParseToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, util.ErrResp("you dont have access"))
		}

		claimsMap := claims.(jwt.MapClaims)
		c.Set("id", uint(claimsMap["id"].(float64)))
		c.Set("role", claimsMap["role"].(string))
		return next(c)
	}
}
