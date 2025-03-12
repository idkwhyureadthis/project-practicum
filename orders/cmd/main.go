package main

import (
	"fmt"
	"os"

	"github.com/idkwhyureadthis/project-practicum/orders/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	connUrl := os.Getenv("DB_URL")
	app := app.New(connUrl)
	fmt.Println(app)
	app.Run(port)
}
