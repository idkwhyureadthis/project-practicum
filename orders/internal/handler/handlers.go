package handler

import (
	"net/http"
	"time"

	"github.com/idkwhyureadthis/project-practicum/orders/internal/service"
	"github.com/labstack/echo/v4"
)

// LogIn godoc
// @Summary Войти в систему
// @Description Аутентификация по номеру телефона и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LogInRequest true "Данные для входа"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *Handler) LogIn(c echo.Context) error {
	var data LogInRequest
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	if data.PhoneNumber == "" {
		data.PhoneNumber = c.QueryParam("phone_number")
	}
	if data.Password == "" {
		data.Password = c.QueryParam("password")
	}

	tokens, user, err := h.s.LogIn(data.PhoneNumber, data.Password)
	if err != nil {
		var code int
		switch err {
		case service.ErrWrongLoginOrPass:
			code = http.StatusUnauthorized
		default:
			code = http.StatusInternalServerError
		}
		return c.JSON(code, map[string]interface{}{
			"error": err.Error(),
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh"
	cookie.Value = tokens.Refresh
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, AuthResponse{
		Tokens: TokensResponse{
			Access:  tokens.Access,
			Refresh: tokens.Refresh,
		},
		User: UserResponse{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
			Mail:        user.Mail,
		},
	})
}

// SignUp godoc
// @Summary Зарегистрироваться
// @Description Создание нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "Данные регистрации"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /signup [post]
func (h *Handler) SignUp(c echo.Context) error {
	var data SignUpRequest
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	user, tokens, code, err := h.s.SignUp(data.PhoneNumber, data.Password, data.Name, data.Email)
	if err != nil {
		return c.JSON(code, ErrorResponse{Error: err.Error()})
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh"
	cookie.Value = tokens.Refresh
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, AuthResponse{
		Tokens: TokensResponse{
			Access:  tokens.Access,
			Refresh: tokens.Refresh,
		},
		User: UserResponse{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
			Mail:        user.Mail,
		},
	})
}

type LogInRequest struct {
	PhoneNumber string `json:"phone_number" example:"+79991234567"`
	Password    string `json:"password" example:"qwerty123"`
}

type SignUpRequest struct {
	Name        string `json:"name" example:"Иван Иванов"`
	Password    string `json:"password" example:"qwerty123"`
	PhoneNumber string `json:"phone_number" example:"+79991234567"`
	Email       string `json:"email" example:"ivan@mail.ru"`
}
