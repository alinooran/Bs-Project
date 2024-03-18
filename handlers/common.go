package handlers

import (
	"github.com/alinooran/Bs-Project/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InternalServerError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, util.ErrResp("خطایی در سرور رخ داده است"))
}
