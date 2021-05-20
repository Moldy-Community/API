package main

import (
	routes "moldy-api/routes"
	"moldy-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	routes.Router(r)
	if utils.GetEnv("DEPLOY") == "on" {
		r.Run()
	} else {
		r.Run("3000")
	}
}
