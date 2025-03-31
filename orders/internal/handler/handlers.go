package handler

import "github.com/labstack/echo/v4"

func (h *Handler) LogIn(c echo.Context) error {
	user, code, err := h.s.LogIn(c.QueryParam("phone_number"), c.QueryParam("password"))
	if err != nil {
		return c.JSON(code, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(code, user)
}

func (h *Handler) SignUp(c echo.Context) error {
	data := struct {
		Name        string `json:"name"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}{}
	if err := c.Bind(&data); err != nil {
		return c.JSON(400, map[string]interface{}{"error": err.Error()})
	}
	user, code, err := h.s.SignUp(data.PhoneNumber, data.Password, data.Name, data.Email)
	if err != nil {
		return c.JSON(code, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(code, user)
}
