package handler

import "github.com/labstack/echo/v4"

func Err2Json(message string, c echo.Context, code int) error {
	return c.JSON(code, map[string]interface{}{
		"error": message,
	})
}
