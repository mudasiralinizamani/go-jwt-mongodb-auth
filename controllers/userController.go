package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mudasiralinizamani/ums-api/helpers"
	"github.com/mudasiralinizamani/ums-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")

		if err := helpers.MatchUserToUid(c, user_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"userid": user_id}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
