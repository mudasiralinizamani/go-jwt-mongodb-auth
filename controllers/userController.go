package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	}
}
