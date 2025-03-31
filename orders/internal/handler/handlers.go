package handler

import (
	"net/http"
	"time"

	"github.com/idkwhyureadthis/project-practicum/orders/internal/service"
	"github.com/labstack/echo/v4"
)

func (h *Handler) LogIn(c echo.Context) error {
	var data struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}
	
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
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
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"tokens": tokens,
		"user": user,
	})
}

func (h *Handler) SignUp(c echo.Context) error {
	data := struct {
		Name        string `json:"name"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}{}
	
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	user, tokens, code, err := h.s.SignUp(data.PhoneNumber, data.Password, data.Name, data.Email)
	if err != nil {
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
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"tokens": tokens,
		"user": user,
	})
}