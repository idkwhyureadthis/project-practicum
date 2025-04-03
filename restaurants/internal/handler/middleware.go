package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearer := c.Request().Header.Get(echo.HeaderAuthorization)
		bearerSplitted := strings.Split(bearer, " ")
		if len(bearerSplitted) != 2 || bearerSplitted[0] != "Bearer" {
			return Err2Json("wrong auth header provided", c, http.StatusUnauthorized)
		}

		token := bearerSplitted[1]
		role, err := h.s.Verify(token, "access")
		if err != nil {
			return Err2Json("wrong auth header provided", c, http.StatusUnauthorized)
		}

		c.Set("role", *role)
		return next(c)
	}
}
