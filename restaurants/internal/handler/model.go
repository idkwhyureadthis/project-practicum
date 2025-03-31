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
	h.e.POST("/verify", h.verify)
	h.e.POST("/generate", h.generate)
}
