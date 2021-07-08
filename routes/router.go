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

	packagesAuth := route.Group("/api/v1/packages", routes.AuthUser)
	{
		packagesAuth.POST("/new", routes.NewPackage)
		packagesAuth.PUT("/update/:id", routes.UpdatePackage)
		packagesAuth.DELETE("/delete/:id", routes.DeleteOne)
	}

	packages := route.Group("/api/v1/packages")
	{
		packages.GET("/search", routes.SearchMany)
		packages.GET("/:id", routes.SearchId)
	}

	users := route.Group("/")
	{
		users.POST("/signup", routes.SignUp)
		users.POST("/login", routes.Login)
	}

	route.GET("/validate-token", routes.ValidToken)

	route.NoRoute(routes.NotFound)
}
