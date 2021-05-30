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
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var packageCollection = database.GetCollection("packages")

func GetAll(c *gin.Context) {
	var allData models.Packages

	cursor, err := packageCollection.Find(context.TODO(), primitive.D{})

	utils.CheckErrors(err, "code 4", "Search finished", "No solution. The search finish")

	for cursor.Next(context.Background()) {
		var pkg models.Format
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
			"data":    nil,
		})
		return
	}

	repeatedName := functions.RepeatedData(reqBody.Name)

	if repeatedName {
		c.JSON(409, gin.H{
			"error":   true,
			"message": "The name of this package was used before",
			"data":    nil,
		})
		return
	}

	re := regexp.MustCompile(`[^0-9|.]`)

	invalid := re.MatchString(reqBody.Version)

	if invalid || strings.HasPrefix(reqBody.Version, ".") {
		c.JSON(406, gin.H{
			"error":   true,
			"message": "Please provide a valid version (Only numbers and dot's)",
			"data":    nil,
		})
		return
	}

	if len(reqBody.Version) > 5 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "The length of the blank <version> is too long",
			"data":    nil,
		})
		return
	} else if len(reqBody.Version) < 3 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "The version is too short, remind that the format is X.Y or X.Y.Z",
			"data":    nil,
		})
		return
	}

	if reqBody.Author == "" || reqBody.Description == "" || reqBody.Name == "" || reqBody.Url == "" || reqBody.Password == "" {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "Please fill all blanks",
			"data":    nil,
		})
		return
	}

	if len(reqBody.Password) < 6 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "Write a password more long (6+ Characters)",
			"data":    nil,
		})
	}

	if len(reqBody.Author) >= 30 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "Please enter a valid author with less of 30 characters",
			"data":    nil,
		})
		return
	}

	if len(reqBody.Description) >= 150 {
		c.JSON(411, gin.H{
			"error":   true,
			"message": "The description is too long and have more of 150 characters, please write it more short",
			"data":    nil,
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
	reqBody.ID = uuid.New().String()

	_, err = packageCollection.InsertOne(context.Background(), reqBody)

	utils.CheckErrors(err, "code 4", "Failed to save in the collection", "Unknown solution")

	formated := &models.Format{
		ID:          reqBody.ID,
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

func UpdatePackage(c *gin.Context) {
	id := c.Param("id")

	var reqBody models.PackageUpdate

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "The request was invalid",
			"data":    nil,
		})
		return
	}

	matchPasswords := functions.SamePassword(reqBody.Password, id)

	if !matchPasswords {
		c.JSON(403, gin.H{
			"error":   true,
			"message": "The password not is correct",
			"data":    nil,
		})
		return
	}

	reqBody.Password = functions.Encrypt(reqBody.Password)

	filter := bson.M{"id": id}

	update := bson.M{
		"$set": bson.M{
			"name":        reqBody.Name,
			"description": reqBody.Description,
			"author":      reqBody.Author,
			"password":    reqBody.Password,
			"url":         reqBody.Url,
			"version":     reqBody.Version,
		},
	}

	_, err := packageCollection.UpdateOne(context.TODO(), filter, update)

	utils.CheckErrors(err, "code 4", "The package failed to be updated", "Try again and watch if the ID is correct")

	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "The package cannot be updated for some unknown reason",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "The package was updated succesfully",
		"data": models.Format{
			ID:          id,
			Name:        reqBody.Name,
			Author:      reqBody.Author,
			Url:         reqBody.Url,
			Description: reqBody.Description,
			Version:     reqBody.Version,
		},
	})
}

func DeleteOne(c *gin.Context) {
	id := c.Param("id")
	var reqBody models.AuthPassword

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "The request was invalid",
			"data":    nil,
		})
		return
	}

	matchPasswords := functions.SamePassword(reqBody.Password, id)

	if !matchPasswords {
		c.JSON(403, gin.H{
			"error":   true,
			"message": "The id or the password is not correct",
			"data":    nil,
		})
		return
	}

	_, err := packageCollection.DeleteOne(context.TODO(), primitive.M{"id": id})

	utils.CheckErrors(err, "code 4", "Something happen in the db and cannot be deleted the package", "Unknown solution")

	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Something bad happen, the package was cannot deleted",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "The package with the id " + id + " was deleted successfully",
		"data":    nil,
	})
}

func SearchId(c *gin.Context) {
	var structure models.Format
	id := c.Param("id")

	err := packageCollection.FindOne(context.TODO(), primitive.D{{Key: "id", Value: id}}).Decode(&structure)

	if err != nil {
		c.JSON(404, gin.H{
			"error":   true,
			"message": "The package not was found by this ID, please verify if is correct.",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Success search by ID",
		"data":    structure,
	})
}

func SearchMany(c *gin.Context) {
	name, found := c.GetQuery("key")

	if !found {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "Bad request, please provide a query param",
			"data":    nil,
		})
		return
	}

	var allData models.Packages

	cursor, err := packageCollection.Find(context.TODO(), bson.M{"name": bson.M{"$regex": `(?i)` + name}})

	utils.CheckErrors(err, "code 4", "Search finished", "No solution. The search finish")

	for cursor.Next(context.Background()) {
		var pkg models.Format
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
