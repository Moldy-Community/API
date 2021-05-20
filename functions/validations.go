package functions

import (
	"context"
	"moldy-api/database"
	"moldy-api/models"
	"moldy-api/utils"
	"strings"

	"github.com/alexedwards/argon2id"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var packageCollection = database.GetCollection("packages")

func RepeatedData(name string) bool {
	var structure = models.Package{}
	err := packageCollection.FindOne(context.TODO(), primitive.D{{Key: "name", Value: name}}).Decode(&structure)

	if err != nil {
		return false
	}

	return strings.EqualFold(structure.Name, name)
}

func SamePassword(password, id string) bool {
	var structure models.Package
	err := packageCollection.FindOne(context.TODO(), primitive.D{{Key: "id", Value: id}}).Decode(&structure)

	utils.CheckErrors(err, "code 4", "The package not was found", "Verify the ID provided")

	match, err := argon2id.ComparePasswordAndHash(password, structure.Password)

	utils.CheckErrors(err, "code 5", "Unknown error validating the password", "Unknown solution")

	return match
}
