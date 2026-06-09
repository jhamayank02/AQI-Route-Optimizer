package router

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {

	rg := r.Group("/api")

	rg.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
}
