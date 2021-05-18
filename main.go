package main

import (
	"log"

	routes "moldy-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.Router(r)

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err.Error())
	}
}
