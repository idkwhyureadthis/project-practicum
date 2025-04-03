package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "needed authorization",
			})
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "wrong token format",
			})
		}

		token := authHeader[7:]

		userID, err := h.s.Verify(token, "access")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": err.Error(),
			})
		}

		user, err := h.s.GetUserByID(userID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "user not found",
			})
		}

		c.Set("user", user)
		c.Set("userID", userID)

		return next(c)
	}
}

func (h *Handler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil {
		data := struct {
			RefreshToken string `json:"refresh_token"`
		}{}

		if err := c.Bind(&data); err != nil || data.RefreshToken == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "refresh token not provided",
			})
		}

		tokens, err := h.s.Refresh(data.RefreshToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, tokens)
	}

	tokens, err := h.s.Refresh(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	newCookie := new(http.Cookie)
	newCookie.Name = "refresh"
	newCookie.Value = tokens.Refresh
	newCookie.Expires = cookie.Expires
	newCookie.HttpOnly = true
	newCookie.Path = "/"
	c.SetCookie(newCookie)

	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) GetProfile(c echo.Context) error {
	user := c.Get("user")
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Logout(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to fetch user",
		})
	}

	err := h.s.InvalidateRefreshToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error logging out",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh"
	cookie.Value = ""
	cookie.Expires = http.Cookie{}.Expires
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success logging out",
	})
}
