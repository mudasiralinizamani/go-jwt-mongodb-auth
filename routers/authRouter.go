package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mudasiralinizamani/ums-api/controllers"
)

func AuthRouter(app *gin.Engine) {
	app.POST("/auth/signup", controllers.Signup())
	app.POST("/auth/signin", controllers.Signin())
}
