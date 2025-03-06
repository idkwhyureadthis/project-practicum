package app

import (
	"log"

	"github.com/idkwhyureadthis/project-practicum/restaurant-service/internal/handler"
)

type App struct {
	h *handler.Handler
}

func New(connUrl string) *App {
	app := App{}
	app.h = handler.New(connUrl)

	return &app
}

func (a *App) Run(port string) {
	log.Fatal(a.h.Run(":" + port))
}
