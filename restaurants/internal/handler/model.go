package handler

import (
	"github.com/idkwhyureadthis/project-practicum/restaurant-service/internal/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	e *echo.Echo
	s *service.Service
}

func (h *Handler) Run(port string) error {
	return h.e.Start(port)
}

func New(connUrl string) *Handler {
	handler := Handler{}
	handler.e = echo.New()
	handler.s = service.New(connUrl)

	return &handler
}
