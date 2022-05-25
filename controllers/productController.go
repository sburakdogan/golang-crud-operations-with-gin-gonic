package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/burak/product-api/database"
	"github.com/burak/product-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productCollection = database.GetCollection(database.GetMongoClient(), "products")

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		defer cancel()

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		validationErrors := validate.Struct(&product)
		if validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": validationErrors.Error()})
			return
		}

		product.Id = primitive.NewObjectID()
		product.CreatedOn = time.Now()

		insertedId, err := productCollection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, insertedId)
	}
}

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var products []models.Product
		defer cancel()

		page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
		limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)

		findOptions := options.Find()

		findOptions.SetSkip((page - 1) * limit)
		findOptions.SetLimit(limit)

		results, err := productCollection.Find(ctx, bson.M{}, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		defer results.Close(ctx)

		for results.Next(ctx) {
			var singleProduct models.Product
			if err = results.Decode(&singleProduct); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			}

			products = append(products, singleProduct)
		}

		c.JSON(http.StatusOK, products)
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		defer cancel()

		productId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		err := productCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		defer cancel()

		productId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		err := productCollection.FindOneAndDelete(ctx, bson.M{"_id": productId}).Decode(&product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": productId})
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		defer cancel()

		productId, _ := primitive.ObjectIDFromHex(c.Param("id"))

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		validationErrors := validate.Struct(&product)
		if validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": validationErrors.Error()})
			return
		}

		update := bson.M{"title": product.Title, "description": product.Description, "categoryName": product.CategoryName}

		_, err := productCollection.UpdateOne(ctx, bson.M{"_id": productId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, update)
	}
}
