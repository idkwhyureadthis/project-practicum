package app

import (
	"log"

	"github.com/idkwhyureadthis/project-practicum/orders/internal/handler"
)

type App struct {
	h *handler.Handler
}

func New(connUrl string, secret string) *App {
	app := App{}
	app.h = handler.New(connUrl, secret)

	return &app
}

func (a *App) Run(port string) {
	log.Fatal(a.h.Run(":" + port))
}
