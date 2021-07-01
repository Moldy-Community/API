package functions

import (
	"context"
	"moldy-api/database"
	"moldy-api/models"
	"strings"

	"github.com/alexedwards/argon2id"
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

func RepeatedUser(name string) bool {

	var structure = models.User{}
	err := userCollection.FindOne(context.TODO(), primitive.D{{Key: "name", Value: name}}).Decode(&structure)

	if err != nil {
		return false
	}

	return strings.EqualFold(structure.Name, name)
}

func SamePassword(password, email string) (bool, string) {
	var structure models.User
	err := userCollection.FindOne(context.TODO(), primitive.D{{Key: "email", Value: email}}).Decode(&structure)

	if err != nil {
		return false, ""
	}

	match, _ := argon2id.ComparePasswordAndHash(password, structure.Password)

	return match, structure.Name
}
