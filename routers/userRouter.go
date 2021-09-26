package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mudasiralinizamani/ums-api/controllers"
)

func UserRouter(app *gin.Engine) {
	app.GET("/users", controllers.GetUsers())
}
