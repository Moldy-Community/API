package routes

import "github.com/gin-gonic/gin"

func GetResponse(c *gin.Context) {
	c.JSON(200, gin.H{
		"error":   false,
		"message": "success",
		"data":    nil,
	})
}

func NotFound(c *gin.Context) {
	c.JSON(404, gin.H{
		"error":   true,
		"message": "Not found route",
		"data":    nil,
	})
}
