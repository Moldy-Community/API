package main

import (
	routes "moldy-api/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.Router(r)
	port, err := os.Getenv("PORT")
	if err != nil {
		port = "3000"
	}
	r.Run(":" + port)
}
