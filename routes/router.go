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
		packages.GET("/search", routes.SearchMany)
		packages.GET("/search/one", routes.SearchOne)
		packages.GET("/:id", routes.SearchId)
		packages.POST("/new", routes.NewPackage)
		packages.PUT("/update/:id", routes.UpdatePackage)
		packages.DELETE("/delete/:id", routes.DeleteOne)
	}

	route.NoRoute(routes.NotFound)
}
