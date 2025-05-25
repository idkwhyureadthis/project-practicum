package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db/generated"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware godoc
// @security BearerAuth
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

// Refresh godoc
// @Summary Обновить токены
// @Description Обновляет access и refresh токены
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body handler.RefreshRequest false "Refresh токен"
// @Success 200 {object} TokensResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /refresh [post]
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

	return c.JSON(http.StatusOK, TokensResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	})
}

// GetProfile godoc
// @Summary Получить профиль
// @Description Возвращает данные авторизованного пользователя
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /profile [get]
func (h *Handler) GetProfile(c echo.Context) error {
	userRaw := c.Get("user")
	dbUser, ok := userRaw.(*generated.GetUserByIDRow)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid user data"})
	}

	return c.JSON(http.StatusOK, dbUser)
}

// Logout godoc
// @Summary Выйти из системы
// @Description Инвалидирует refresh токен
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} LogoutResponse
// @Failure 500 {object} ErrorResponse
// @Router /logout [post]
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

	return c.JSON(http.StatusOK, LogoutResponse{Message: "success logging out"})
}
