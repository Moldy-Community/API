package database

import (
	"context"
	utils "moldy-api/utils"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(collection string) *mongo.Collection {
	URI := utils.GetEnv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))

	utils.CheckErrors(err, "1", "The connection to the database failed", "Please check if the URI value in the .env is valid")

	database := client.Database("community-api")

	return database.Collection(collection)
}
