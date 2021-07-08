package routes

import (
	"context"
	"fmt"
	"moldy-api/database"
	"moldy-api/functions"
	"moldy-api/models"
	"moldy-api/utils"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var packageCollection = database.GetCollection("packages")

func SearchId(c *gin.Context) {
	var structure models.Format
	id := c.Param("id")
	token := c.Request.Header.Get("token")

	newToken := functions.NewToken(token)

	err := packageCollection.FindOne(context.TODO(), primitive.D{{Key: "id", Value: id}}).Decode(&structure)

	if err != nil {
		c.JSON(404, gin.H{
			"error":    true,
			"message":  "The package not was found by this ID, please verify if is correct.",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":    false,
		"message":  "Success search by ID",
		"data":     structure,
		"newToken": newToken,
	})
}

func SearchMany(c *gin.Context) {
	name, found := c.GetQuery("query")
	limitStr, _ := c.GetQuery("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if limit == 0 || err != nil {
		limit = 20
	}
	token := c.Request.Header.Get("token")

	newToken := functions.NewToken(token)

	if !found {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "Bad request, please provide a query param",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	var allData models.Packages

	cursor, err := packageCollection.Find(context.TODO(), bson.M{"name": bson.M{"$regex": `(?i)` + name}})

	utils.CheckErrors(err, "code 4", "Search finished", "No solution. The search finish")

	var i int64

	for i = 0; cursor.Next(context.Background()); i++ {
		if i < limit {
			var pkg models.Format
			_ = cursor.Decode(&pkg)

			allData = append(allData, &pkg)

		} else {
			break
		}
	}

	c.JSON(200, gin.H{
		"error":    false,
		"message":  "Success search",
		"data":     allData,
		"newToken": newToken,
	})

}

func NewPackage(c *gin.Context) {
	var reqBody models.Package
	token := c.Request.Header.Get("token")

	newToken := functions.NewToken(token)

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "The request was invalid",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	repeatedName := functions.RepeatedPackage(reqBody.Name)

	if repeatedName {
		c.JSON(409, gin.H{
			"error":    true,
			"message":  "The name of this package was used before",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	re := regexp.MustCompile(`[^0-9|.]`)

	invalid := re.MatchString(reqBody.Version)

	if invalid || strings.HasPrefix(reqBody.Version, ".") {
		c.JSON(406, gin.H{
			"error":    true,
			"message":  "Please provide a valid version (Only numbers and dot's)",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	if len(reqBody.Version) > 5 {
		c.JSON(411, gin.H{
			"error":    true,
			"message":  "The length of the blank <version> is too long",
			"data":     nil,
			"newToken": newToken,
		})
		return
	} else if len(reqBody.Version) < 3 {
		c.JSON(411, gin.H{
			"error":    true,
			"message":  "The version is too short, remind that the format is X.Y or X.Y.Z",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	if reqBody.Description == "" || reqBody.Name == "" || reqBody.Url == "" {
		c.JSON(411, gin.H{
			"error":    true,
			"message":  "Please fill all blanks",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	if len(reqBody.Description) >= 150 {
		c.JSON(411, gin.H{
			"error":    true,
			"message":  "The description is too long and have more of 150 characters, please write it more short",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	u, err := url.Parse(reqBody.Url)

	if err != nil || u.Host == "" || u.Scheme == "" {
		c.JSON(406, gin.H{
			"error":    true,
			"message":  "The URL is not valid",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	reqBody.ID = uuid.New().String()

	validatedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(utils.GetEnv("JWT_SIGN")), nil
	})
	claims, _ := validatedToken.Claims.(jwt.MapClaims)
	reqBody.Author = string(claims["name"].(string))

	_, err = packageCollection.InsertOne(context.Background(), reqBody)

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happened when the data was saved",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	formated := &models.Format{
		ID:          reqBody.ID,
		Name:        reqBody.Name,
		Author:      reqBody.Author,
		Url:         reqBody.Url,
		Description: reqBody.Description,
		Version:     reqBody.Version,
	}

	c.JSON(200, gin.H{
		"error":    false,
		"message":  "Created successfully",
		"data":     formated,
		"newToken": newToken,
	})

}

func UpdatePackage(c *gin.Context) {
	id := c.Param("id")
	var reqBody models.PackageUpdate
	token := c.Request.Header.Get("token")

	newToken := functions.NewToken(token)

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "The request was invalid",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	validatedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(utils.GetEnv("JWT_SIGN")), nil
	})
	claims, _ := validatedToken.Claims.(jwt.MapClaims)
	reqBody.Author = string(claims["name"].(string))

	filter := bson.M{"id": id}

	update := bson.M{
		"$set": bson.M{
			"name":        reqBody.Name,
			"description": reqBody.Description,
			"author":      reqBody.Author,
			"url":         reqBody.Url,
			"version":     reqBody.Version,
		},
	}

	_, err := packageCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happened when the data was attempted to save",
			"data":     nil,
			"newToken": newToken,
		})
	}

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "The package cannot be updated for some unknown reason",
			"data":     nil,
			"newToken": newToken,
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
		"newToken": newToken,
	})
}

func DeleteOne(c *gin.Context) {
	var structure models.Package
	id := c.Param("id")
	token := c.Request.Header.Get("token")

	newToken := functions.NewToken(token)

	err := packageCollection.FindOne(context.TODO(), primitive.D{{Key: "id", Value: id}}).Decode(&structure)

	validatedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(utils.GetEnv("JWT_SIGN")), nil
	})
	claims, _ := validatedToken.Claims.(jwt.MapClaims)

	if structure.Author != string(claims["name"].(string)) {
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "Unauthorized",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	_, err = packageCollection.DeleteOne(context.TODO(), primitive.M{"id": id})

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happen, the package was cannot deleted",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":    false,
		"message":  "The package with the id " + id + " was deleted successfully",
		"data":     nil,
		"newToken": newToken,
	})
}
