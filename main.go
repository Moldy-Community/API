package main

import (
	routes "moldy-api/routes"

	"github.com/gin-gonic/gin"
)

var deploy bool = true

func main() {
	r := gin.Default()

	routes.Router(r)
	if deploy == true {
		r.Run()
	} else {
		r.Run("3000")
	}
}
