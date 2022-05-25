package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/burak/product-api/database"
	"github.com/burak/product-api/helpers"
	"github.com/burak/product-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = database.GetCollection(database.GetMongoClient(), "users")
var validate = validator.New()

type LoginRequest struct {
	Email    string
	Password string
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		validationErrors := validate.Struct(&user)
		if validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": validationErrors.Error()})
			return
		}

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		if emailCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"err": "This email address is already exist"})
			return
		}

		user.Id = primitive.NewObjectID()
		user.CreatedOn = time.Now()
		user.Password = helpers.HashPassword(user.Password)

		instertedId, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, instertedId)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()

		var loginRequest LoginRequest
		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Login failed. Check your inputs."})
			return
		}

		err = helpers.VerifyPassword(loginRequest.Password, user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Login failed. Check your inputs."})
			return
		}

		token, err := helpers.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": user.Id, "fullName": user.FullName, "email": user.Email, "token": token})
	}
}
