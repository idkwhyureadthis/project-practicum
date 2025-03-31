package handler

import (
	"net/http"
	"time"

	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/service"
	"github.com/labstack/echo/v4"
)

func (h *Handler) login(c echo.Context) error {
	var data LoginRequest

	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}
	tokens, err := h.s.LogIn(data.Login, data.Password)
	if err == service.ErrWrongData {
		return Err2Json(err.Error(), c, http.StatusUnauthorized)
	} else if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh"
	cookie.Value = tokens.Refresh
	cookie.Expires = time.Now().Add(time.Hour * 24)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) verify(c echo.Context) error {
	var data VerifyRequest
	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	role, err := h.s.Verify(data.Token, "access")
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"role": &role,
	})
}

func (h *Handler) generate(c echo.Context) error {
	var data GenerateRequest
	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	tokens, err := h.s.Generate(data.Token)
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, tokens)
}
