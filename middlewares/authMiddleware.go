package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mudasiralinizamani/go-jwt-mongodb-auth/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("uid", claims.Uid)
		c.Set("role", claims.Role)
		c.Next()
	}
}
