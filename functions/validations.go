package functions

import (
	"context"
	"fmt"
	"moldy-api/database"
	"moldy-api/models"
	"moldy-api/utils"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var packageCollection = database.GetCollection("packages")
var userCollection = database.GetCollection("users")

func RepeatedPackage(name string) bool {
	var structure = models.Package{}
	err := packageCollection.FindOne(context.TODO(), primitive.D{{Key: "name", Value: name}}).Decode(&structure)

	if err != nil {
		return false
	}

	return strings.EqualFold(structure.Name, name)
}

func RepeatedUser(value, key string) bool {
	var structure = models.User{}
	err := userCollection.FindOne(context.TODO(), primitive.D{{Key: key, Value: value}}).Decode(&structure)

	if err != nil {
		return false
	}

	if key == "email" {
		return strings.EqualFold(structure.Email, value)
	} else if key == "name" {
		return strings.EqualFold(structure.Name, value)
	}

	return true
}

func SamePassword(password, email string) (bool, *models.User) {
	var structure *models.User
	err := userCollection.FindOne(context.TODO(), primitive.D{{Key: "email", Value: email}}).Decode(&structure)

	if err != nil {
		return false, &models.User{}
	}

	match, _ := argon2id.ComparePasswordAndHash(password, structure.Password)

	return match, structure
}

func TokenIsValid(token string) (bool, jwt.MapClaims) {
	validatedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(utils.GetEnv("JWT_SIGN")), nil
	})

	claims, ok := validatedToken.Claims.(jwt.MapClaims)

	if ok && validatedToken.Valid {
		return true, claims
	}

	return false, nil
}

func NewToken(token string) interface{} {
	if len(token) != 0 {
		isValid, claims := TokenIsValid(token)
		if isValid {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email":    claims["email"],
				"name":     claims["name"],
				"expireAt": time.Now().Format(time.RFC3339),
			})

			signedToken, err := token.SignedString([]byte(utils.GetEnv("JWT_SIGN")))

			if err != nil {
				return nil
			}

			return signedToken
		}
	}
	return nil
}
