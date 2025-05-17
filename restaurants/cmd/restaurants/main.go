package main

import (
	"os"

	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	connUrl := os.Getenv("DB_URL")
	secretKey := os.Getenv("JWT_SECRET")
	app := app.New(connUrl, adminPass, secretKey)
	app.Run(port)
}
