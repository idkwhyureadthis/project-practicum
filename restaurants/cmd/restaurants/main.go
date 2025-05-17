package main

import (
	"os"

	_ "github.com/idkwhyureadthis/project-practicum/restaurants/docs"
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/app"
	"github.com/joho/godotenv"
)

// @title           Restaurants API
// @version         1.0
// @description     API for managing restaurants, admins and items.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	connUrl := os.Getenv("DB_URL")
	adminPass := os.Getenv("ADMIN_PASS")
	secretKey := os.Getenv("JWT_SECRET")
	app := app.New(connUrl, adminPass, secretKey)
	app.Run(port)
}
