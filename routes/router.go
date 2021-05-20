package routes

import (
	routes "moldy-api/controllers/routes-controller"

	"github.com/gin-gonic/gin"
)

func Router(route *gin.Engine) {
	main := route.Group("/")
	{
		main.GET("/", routes.GetResponse)
	}

	packages := route.Group("/api/v1/packages")
	{
		packages.GET("/all", routes.GetAll)
		packages.POST("/new", routes.NewPackage)
	}

	route.NoRoute(routes.NotFound)
}
