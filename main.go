package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mudasiralinizamani/ums-api/routers"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	app := gin.New()
	app.Use(gin.Logger())

	// Adding Router
	routers.AuthRouter(app)
	routers.UserRouter(app)

	app.Run(":" + port)
}
