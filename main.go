package main

import (
	routes "moldy-api/routes"
	"moldy-api/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var deploy bool = true

func main() {
	err := godotenv.Load()

	utils.CheckErrors(err, "2", "Error in read the enviroment variables")

	r := gin.Default()

	routes.Router(r)
	if os.Getenv("DEPLOY") == "on" {
		r.Run()
	} else {
		r.Run("3000")
	}
}
