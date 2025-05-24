package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/service"
	"github.com/jackc/pgx/v5"
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

func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	uid, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("userID not found in context")
	}
	return uid, nil
}

// CreateOrder godoc
// @Summary Создать заказ
// @Description Создание нового заказа в ресторане
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "Параметры заказа"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (h *Handler) CreateOrder(c echo.Context) error {
	var req CreateOrderRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
	}

	order, _ := h.s.CreateOrder(req.DisplayedID, &req.RestaurantID, req.TotalPrice, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	resp := OrderResponse{
		ID:           order.ID,
		DisplayedID:  order.DisplayedID,
		RestaurantID: *order.RestaurantID,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
	}

	return c.JSON(http.StatusOK, resp)
}

// GetOrderByID получает заказ по UUID
// @Summary Получение заказа
// @Tags Orders
// @Security BearerAuth
// @Description Получает заказ по его UUID
// @Param id path string true "UUID заказа"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} ErrorResponse "Некорректный UUID"
// @Failure 404 {object} ErrorResponse "Заказ не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/{id} [get]
func (h *Handler) GetOrderByID(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
	}

	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Некорректный UUID заказа"})
	}

	order, err := h.s.GetOrderByID(orderID, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "Заказ не найден"})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	resp := OrderResponse{
		ID:           order.ID,
		DisplayedID:  order.DisplayedID,
		RestaurantID: *order.RestaurantID,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
	}
	return c.JSON(http.StatusOK, resp)
}

// GetAllOrders получает все заказы
// @Summary Получение всех заказов
// @Tags Orders
// @Security BearerAuth
// @Description Получает список всех заказов (без фильтрации по пользователю)
// @Success 200 {array} OrderResponse
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders [get]
func (h *Handler) GetAllOrders(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
	}

	orders, err := h.s.GetAllOrders(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	var resp []OrderResponse
	for _, o := range orders {
		resp = append(resp, OrderResponse{
			ID:           o.ID,
			DisplayedID:  o.DisplayedID,
			RestaurantID: *o.RestaurantID,
			TotalPrice:   o.TotalPrice,
			Status:       o.Status,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteOrder удаляет заказ по UUID
// @Summary Удаление заказа
// @Tags Orders
// @Security BearerAuth
// @Description Удаляет заказ по его UUID
// @Param id path string true "UUID заказа"
// @Success 204 "Заказ удалён"
// @Failure 400 {object} ErrorResponse "Некорректный UUID"
// @Failure 404 {object} ErrorResponse "Заказ не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/{id} [delete]
func (h *Handler) DeleteOrder(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid order id"})
	}

	err = h.s.DeleteOrder(id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "order not found"})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
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

type CreateOrderRequest struct {
	DisplayedID  int32     `json:"displayed_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	TotalPrice   float64   `json:"total_price"`
}
