package main

import (
	routes "moldy-api/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.Router(r)
	port := os.Getenv("PORT")
	port = "3000"
	r.Run(":" + port)
}
