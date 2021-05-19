package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()

	CheckErrors(err, "code 6", "Error loading the .env file", "Create a .env file and write there the value, something like this:\nVALUE_EXAMPLE=ThisIsTheValue")

	value := os.Getenv(key)

	return value
}
