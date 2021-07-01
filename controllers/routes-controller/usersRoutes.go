package routes

import (
	"context"
	"fmt"
	"moldy-api/database"
	"moldy-api/functions"
	"moldy-api/models"
	"moldy-api/utils"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var userCollection = database.GetCollection("users")

func SignUp(c *gin.Context) {
	var reqBody models.User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "The request was invalid",
			"data":    nil,
		})
		return
	}

	if reqBody.Name == "" || reqBody.Email == "" || reqBody.Password == "" {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please fill all blanks",
			"data":    nil,
		})
		return
	}

	repeatedUser := functions.RepeatedUser(reqBody.Name)

	if repeatedUser {
		c.JSON(403, gin.H{
			"error":   true,
			"message": "The user was created before",
			"data":    nil,
		})
		return
	}

	if len(reqBody.Password) < 6 {
		c.JSON(406, gin.H{
			"error":   true,
			"message": "Please write a password with 6 chars or more",
			"data":    nil,
		})
		return
	}

	_, err := mail.ParseAddress(reqBody.Email)

	if err != nil {
		c.JSON(406, gin.H{
			"error":   true,
			"message": "Invalid email",
			"data":    nil,
		})
		return
	}

	reqBody.Password = functions.Encrypt(reqBody.Password)

	_, err = userCollection.InsertOne(context.Background(), reqBody)

	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Something bad happenned when the data was saved",
			"data":    nil,
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": reqBody.Email,
		"name":  reqBody.Name,
	})

	signedToken, err := token.SignedString([]byte(utils.GetEnv("JWT_SIGN")))

	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Something bad happened in the auth",
			"data":    nil,
		})
		return
	}

	c.JSON(201, gin.H{
		"error":   false,
		"message": "The user was created successfully",
		"data":    signedToken,
	})
}

func Login(c *gin.Context) {
	var reqBody models.Login
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "The request was invalid",
			"data":    nil,
		})
		return
	}

	passwordMatch, name := functions.SamePassword(reqBody.Password, reqBody.Email)

	if !passwordMatch {
		c.JSON(403, gin.H{
			"error":   true,
			"message": "The password is incorrect",
			"data":    nil,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": reqBody.Email,
		"name":  name,
	})

	signedToken, err := token.SignedString([]byte(utils.GetEnv("JWT_SIGN")))

	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Something bad happened in the auth",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Success login",
		"data":    signedToken,
	})
}

func AuthUser(c *gin.Context) {
	token := c.Request.Header.Get("token")

	if len(token) == 0 {
		c.Abort()
		c.JSON(403, gin.H{
			"error":   true,
			"message": "Missing the token in the headers",
			"data":    nil,
		})
		return
	}

	validatedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(utils.GetEnv("JWT_SIGN")), nil
	})

	if _, ok := validatedToken.Claims.(jwt.MapClaims); ok && validatedToken.Valid {
		c.Next()
	} else {
		c.Abort()
		c.JSON(403, gin.H{
			"error":   true,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}
}
