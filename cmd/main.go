package main

import (
	"github.com/gin-gonic/gin"
	"ProductsAPI/internal/handlers"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Println("⚠️  No .env file found (using system env)")
    }

	fmt.Println("USE_DYNAMO:", os.Getenv("USE_DYNAMO"))

	router := gin.Default()
	
	//health check for ECS
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
	router.Run(":8080")
}
