package handler

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/labstack/echo/v4"
// )

// func (h *Handler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		bearer := c.Request().Header.Get(echo.HeaderAuthorization)
// 		bearerSplitted := strings.Split(bearer, " ")
// 		if len(bearerSplitted) != 1 {
// 			return c.NoContent(http.StatusUnauthorized)
// 		}
// 		token := bearerSplitted[1]

// 		if err := h.s.VerifyToken(token)

// 		return next(c)
// 	}
// }
