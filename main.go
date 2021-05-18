package main

import (
	routes "moldy-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.Router(r)
	r.Run()
}
