package app

import (
	"log"

	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/handler"
)

type App struct {
	h *handler.Handler
}

func New(connUrl, secret string) *App {
	app := App{}
	app.h = handler.New(connUrl, secret)

	return &app
}

func (a *App) Run(port string) {
	log.Fatal(a.h.Run(":" + port))
}
