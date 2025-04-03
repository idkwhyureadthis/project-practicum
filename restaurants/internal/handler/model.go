package handler

import (
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/service"
	"github.com/labstack/echo/v4"
)

type VerifyRequest struct {
	Token string `json:"access"`
}

type GenerateRequest struct {
	Token string `json:"refresh"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateAdminRequest struct {
	Login        string `json:"login"`
	Password     string `json:"password"`
	RestaurantId string `json:"restaurant_id"`
}

type AddRestaurantRequest struct {
	OpenTime  string  `json:"open_time"`
	CloseTime string  `json:"close_time"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Name      string  `json:"name"`
}

type Handler struct {
	e      *echo.Echo
	s      *service.Service
	Secret []byte
}

func (h *Handler) Run(port string) error {
	return h.e.Start(port)
}

func New(connUrl, adminPass, secret string) *Handler {
	handler := Handler{}
	handler.e = echo.New()
	handler.s = service.New(connUrl, adminPass, secret)
	handler.setup()
	return &handler
}

func (h *Handler) setup() {
	h.e.POST("/login", h.login)
	h.e.POST("/refresh", h.refresh)
	h.e.POST("/verify", h.verify)
	restaurants := h.e.Group("/restaurants")
	restaurants.POST("", h.addRestaurant, h.authMiddleware)
	restaurants.GET("", h.getRestaurants)
	admins := h.e.Group("/admins")
	admins.POST("", h.createAdmin, h.authMiddleware)
	items := h.e.Group("/items")
	items.POST("", h.addItem, h.authMiddleware)
}
