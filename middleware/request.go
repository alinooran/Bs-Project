package middleware

import (
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func ParseRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, util.ErrResp("شناسه درخواست معتبر نیست"))
		}
		c.Set("request_id", uint(reqID))
		return next(c)
	}
}
