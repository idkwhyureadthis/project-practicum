package main

import (
	"fmt"
	"os"

	_ "github.com/idkwhyureadthis/project-practicum/orders/docs"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	connUrl := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	app := app.New(connUrl, jwtSecret)
	fmt.Println(app)
	app.Run(port)
}
