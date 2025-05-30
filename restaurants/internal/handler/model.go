package handler

import (
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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

type ItemActionRequest struct {
	ItemId string `json:"item_id"`
}

type CreateOrderRequest struct {
	ItemsIDs     []string `json:"item_ids"`
	RestaurantID string   `json:"restaurant_id"`
}

type Handler struct {
	e      *echo.Echo
	s      *service.Service
	Secret []byte
}

func (h *Handler) Run(port string) error {
	return h.e.Start(port)
}

func New(connUrl, secret string) *Handler {
	handler := Handler{}
	handler.e = echo.New()
	handler.s = service.New(connUrl, secret)
	handler.setup()

	handler.e.Use(middleware.Logger())
	handler.e.Use(middleware.Recover())
	handler.e.Use(middleware.CORS())

	return &handler
}

func (h *Handler) setup() {
	h.e.GET("/swagger/*", echoSwagger.WrapHandler)
	h.e.POST("/login", h.login)
	h.e.POST("/refresh", h.refresh)
	h.e.POST("/verify", h.verify)
	items := h.e.Group("/items")
	admins := h.e.Group("/admins")
	restaurants := h.e.Group("/restaurants")
	orders := h.e.Group("/orders")
	orders.GET("", h.getOrders, h.authMiddleware)
	orders.POST("", h.createOrder)
	restaurants.GET("", h.getRestaurants)
	restaurants.POST("", h.addRestaurant, h.authMiddleware)
	admins.POST("", h.createAdmin, h.authMiddleware)
	items.GET("", h.getItems)
	items.POST("", h.addItem, h.authMiddleware)
	items.POST("/ban", h.banItem, h.authMiddleware)
	items.POST("/unban", h.unbanItem, h.authMiddleware)
}
