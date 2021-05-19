package routes

import (
	"context"
	"moldy-api/database"
	"moldy-api/models"
	"moldy-api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var packageCollection = database.GetCollection("packages")

func GetAll(c *gin.Context) {
	var allData models.Packages

	cursor, err := packageCollection.Find(context.TODO(), primitive.D{})

	utils.CheckErrors(err, "code 4", "Search finished", "No solution. The search finish")

	for cursor.Next(context.Background()) {
		var pkg models.Package
		err = cursor.Decode(&pkg)

		utils.CheckErrors(err, "code 4", "Search finished", "No solution. The search finish")

		allData = append(allData, &pkg)
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Success search",
		"data":    allData,
	})

}

func NewPackage(c *gin.Context) {
	var reqBody models.Package

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "The request was invalid",
		})
		return
	}

	_, err := packageCollection.InsertOne(context.Background(), reqBody)

	utils.CheckErrors(err, "code 4", "Failed to save in the collection", "Unknown solution")

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Created successfully",
		"data":    reqBody,
	})

}
