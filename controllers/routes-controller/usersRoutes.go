package routes

import (
	"context"
	"fmt"
	"moldy-api/database"
	"moldy-api/functions"
	"moldy-api/models"
	"moldy-api/utils"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = database.GetCollection("users")

const (
	month = (time.Hour * 24) * 30
)

func SignUp(c *gin.Context) {
	var reqBody models.User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "The request was invalid",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	if reqBody.Name == "" || reqBody.Email == "" || reqBody.Password == "" {
		c.JSON(400, gin.H{
			"error":    true,
			"message":  "Please fill all blanks",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	repeatedUser := functions.RepeatedUser(reqBody.Name, "name")
	repeatedEmail := functions.RepeatedUser(reqBody.Email, "email")

	if repeatedUser {
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "The name was used before",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	if repeatedEmail {
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "The email was used before",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	if len(reqBody.Password) < 6 {
		c.JSON(406, gin.H{
			"error":    true,
			"message":  "Please write a password with 6 chars or more",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	_, err := mail.ParseAddress(reqBody.Email)

	if err != nil {
		c.JSON(406, gin.H{
			"error":    true,
			"message":  "Invalid email",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	reqBody.Password = functions.Encrypt(reqBody.Password)

	_, err = userCollection.InsertOne(context.Background(), reqBody)

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happenned when the data was saved",
			"data":     nil,
			"newToken": nil,
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    reqBody.Email,
		"name":     reqBody.Name,
		"expireAt": time.Now().Format(time.RFC3339),
	})

	signedToken, err := token.SignedString([]byte(utils.GetEnv("JWT_SIGN")))

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happened in the auth",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	c.JSON(201, gin.H{
		"error":    false,
		"message":  "The user was created successfully",
		"data":     signedToken,
		"newToken": signedToken,
	})
}

func Login(c *gin.Context) {
	var reqBody models.Login
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "The request was invalid",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	passwordMatch, dataStruct := functions.SamePassword(reqBody.Password, reqBody.Email)

	if !passwordMatch {
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "The password is incorrect",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    reqBody.Email,
		"name":     dataStruct.Name,
		"expireAt": time.Now().Format(time.RFC3339),
	})

	signedToken, err := token.SignedString([]byte(utils.GetEnv("JWT_SIGN")))

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happened in the auth",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":    false,
		"message":  "Success login",
		"data":     signedToken,
		"newToken": signedToken,
	})
}

func UpdateUser(c *gin.Context) {
	var reqBody models.UpdateUser
	token := c.Request.Header.Get("token")

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "The request was invalid",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	passwordMatch, dataStruct := functions.SamePassword(reqBody.OldPassword, reqBody.Email)

	if !passwordMatch {
		newToken := functions.NewToken(token)
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "The password is incorrect",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	filter := bson.M{"email": dataStruct.Email}

	update := bson.M{
		"$set": bson.M{
			"name":     reqBody.Name,
			"email":    reqBody.Email,
			"password": reqBody.NewPassword,
		},
	}
	_, err := userCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		newToken := functions.NewToken(token)
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happened when the data was attempted to save",
			"data":     nil,
			"newToken": newToken,
		})
		return
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    reqBody.Email,
		"name":     reqBody.Name,
		"expireAt": time.Now().Format(time.RFC3339),
	})

	signedToken, err := newToken.SignedString([]byte(utils.GetEnv("JWT_SIGN")))

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happened in the auth",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":    true,
		"message":  "The user was updated successfully",
		"data":     nil,
		"newToken": signedToken,
	})
}

func DeleteUser(c *gin.Context) {
	var reqBody struct{ Password string }
	token := c.Request.Header.Get("token")

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "The request was invalid",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	isValid, claims := functions.TokenIsValid(token)

	if !isValid {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "Invalid token.",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	passwordMatch, dataStruct := functions.SamePassword(reqBody.Password, claims["email"].(string))

	if !passwordMatch {
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "The password not match.",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	_, err := userCollection.DeleteOne(context.TODO(), primitive.M{"email": dataStruct.Email})

	if err != nil {
		c.JSON(500, gin.H{
			"error":    true,
			"message":  "Something bad happened at the moment when the user was attempted to delete",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":    false,
		"message":  "The user was deleted successfully",
		"data":     nil,
		"newToken": nil,
	})

}

func AuthUser(c *gin.Context) {
	token := c.Request.Header.Get("token")

	if len(token) == 0 {
		c.Abort()
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "Missing the token in the headers",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	validatedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(utils.GetEnv("JWT_SIGN")), nil
	})

	if err != nil {
		c.Abort()
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "Unauthorized, invalid token.",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	claims, ok := validatedToken.Claims.(jwt.MapClaims)
	expireAt, _ := time.Parse(time.RFC3339, string(claims["expireAt"].(string)))
	timeElapsed := time.Since(expireAt)

	if ok && validatedToken.Valid && timeElapsed < month {
		c.Next()
	} else {
		c.Abort()
		c.JSON(403, gin.H{
			"error":    true,
			"message":  "Unauthorized, invalid token.",
			"data":     nil,
			"newToken": nil,
		})
		return
	}
}

func ValidToken(c *gin.Context) {
	token := c.Request.Header.Get("token")

	if len(token) == 0 {
		c.JSON(422, gin.H{
			"error":    true,
			"message":  "Provide the token in the headers.",
			"data":     nil,
			"newToken": nil,
		})
		return
	}

	isValid, _ := functions.TokenIsValid(token)

	if isValid {
		c.JSON(200, gin.H{
			"error":    false,
			"message":  "The token is valid.",
			"data":     isValid,
			"newToken": nil,
		})
		return
	}

	if !isValid {
		c.JSON(200, gin.H{
			"error":    false,
			"message":  "The token is invalid.",
			"data":     isValid,
			"newToken": nil,
		})
		return
	}
}
