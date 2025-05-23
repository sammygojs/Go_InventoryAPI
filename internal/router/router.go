package router

import (
	"github.com/gin-gonic/gin"
	"ProductsAPI/internal/handlers"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode) // helpful for test logs
	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		products := api.Group("/products")
		{
			products.GET("", handlers.GetProducts)
			products.GET("/:productID", handlers.GetProduct)
		}
	}
	return router
}
