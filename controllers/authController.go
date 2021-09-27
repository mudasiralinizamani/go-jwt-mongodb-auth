package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mudasiralinizamani/ums-api/database"
	"github.com/mudasiralinizamani/ums-api/helpers"
	"github.com/mudasiralinizamani/ums-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while cheking the email"})
		}

		password := helpers.HashPassword(*user.Password)

		user.Password = &password

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email address already exists"})
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		token, refreshToken, err := helpers.GenerateTokens(*user.Email, *user.FirstName, *user.LastName, *user.Role, user.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while creating tokens"})
			log.Panic(err)
		}

		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while creating user"})
			log.Panic(insertErr)
		}

		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email does not exists"})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(*foundUser.Password, *user.Password)

		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		token, refreshToken, err := helpers.GenerateTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.Role, foundUser.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while creating token and refresh token"})
			log.Panic(err)
			return
		}

		helpers.UpdateTokens(token, refreshToken, foundUser.UserId)

		err = userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserId}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}
