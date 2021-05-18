package routes

import (
	mainController "moldy-api/controllers/main-routes"

	"github.com/gin-gonic/gin"
)

func Router(route *gin.Engine) {
	main := route.Group("/")
	{
		main.GET("/", mainController.GetResponse)
	}
}
