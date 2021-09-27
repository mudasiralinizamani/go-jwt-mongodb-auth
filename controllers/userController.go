package controllers

import (
	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// user_id := c.Param("user_id")

	}
}
