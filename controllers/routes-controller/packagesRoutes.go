package routes

import (
	"context"
	"moldy-api/database"
	"moldy-api/functions"
	"moldy-api/models"
	"moldy-api/utils"
	"net/url"
	"regexp"
	"strings"

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

	re := regexp.MustCompile(`[^0-9|.]`)

	invalid := re.MatchString(reqBody.Version)

	if invalid || strings.HasPrefix(reqBody.Version, ".") {
		c.JSON(406, gin.H{
			"error":   true,
			"message": "Please provide a valid version (Only numbers and dot's)",
		})
		return
	}

	if len(reqBody.Version) > 5 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "The length is too long",
		})
		return
	} else if len(reqBody.Version) < 3 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "The version is too short, remind that the format is X.Y or X.Y.Z",
		})
		return
	}

	if reqBody.Author == "" || reqBody.Description == "" || reqBody.Name == "" || reqBody.Url == "" || reqBody.Password == "" {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "Please fill all blanks",
		})
		return
	}

	if len(reqBody.Password) < 6 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "Write a password more long (6+ Characters)",
		})
	}

	if len(reqBody.Author) >= 30 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "Please enter a valid author with less of 30 characters",
		})
		return
	}

	if len(reqBody.Description) >= 150 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "The description is too long and have more of 150 characters, please write it more short",
		})
		return
	}

	u, err := url.Parse(reqBody.Url)

	if err != nil || u.Host == "" || u.Scheme == "" {
		c.JSON(406, gin.H{
			"error":   true,
			"message": "The URL is not valid",
		})
		return
	}

	reqBody.Password = functions.Encrypt(reqBody.Password)

	_, err = packageCollection.InsertOne(context.Background(), reqBody)

	utils.CheckErrors(err, "code 4", "Failed to save in the collection", "Unknown solution")

	formated := &models.Format{
		Name:        reqBody.Name,
		Author:      reqBody.Author,
		Url:         reqBody.Url,
		Description: reqBody.Description,
		Version:     reqBody.Version,
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Created successfully",
		"data":    formated,
	})

}
