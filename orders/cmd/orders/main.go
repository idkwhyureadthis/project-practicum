package main

import (
	"fmt"
	"os"

	_ "github.com/idkwhyureadthis/project-practicum/orders/docs"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/app"
	"github.com/joho/godotenv"
)

// @title Restaurant & Orders API
// @version 1.0
// @description API для управления ресторанами и заказами
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@project.ru
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	connUrl := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	app := app.New(connUrl, jwtSecret)
	fmt.Println(app)
	app.Run(port)
}
