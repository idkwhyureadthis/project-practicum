package handler

import (
	"github.com/idkwhyureadthis/project-practicum/orders/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	e *echo.Echo
	s *service.Service
	secret []byte
}

func (h *Handler) Run(port string) error {
	return h.e.Start(port)
}

func New(connUrl string, secret string) *Handler {
	handler := Handler{}
	handler.e = echo.New()
	handler.s = service.New(connUrl, secret)
	handler.secret = []byte(secret)

	handler.e.Use(middleware.Logger())
	handler.e.Use(middleware.Recover())
	handler.e.Use(middleware.CORS())
	
	handler.SetupHandlers()
	return &handler
}

func (h *Handler) SetupHandlers() {
	h.e.POST("/login", h.LogIn)
	h.e.POST("/signup", h.SignUp) 
	h.e.POST("/refresh", h.Refresh)

	auth := h.e.Group("")
	auth.Use(h.AuthMiddleware)

	auth.GET("/profile", h.GetProfile)
	auth.POST("/logout", h.Logout)
}