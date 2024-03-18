package util

import "github.com/labstack/echo/v4"

func ErrResp(msg string) echo.Map {
	return echo.Map{
		"error": echo.Map{
			"message": msg,
		},
	}
}