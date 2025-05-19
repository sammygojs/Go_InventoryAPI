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

	// Just for confirmation
	fmt.Println("USE_DYNAMO:", os.Getenv("USE_DYNAMO"))

	router := gin.Default()
	
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	api := router.Group("/api")
	{
		fmt.Println("✅ Registered /api/query-db route")
		api.GET("/query-db", handlers.QueryDB) // ✅ Correct group

		products := api.Group("/products")
		{
			products.GET("/", handlers.GetProducts)
			products.GET("/:productID", handlers.GetProduct)
		}
	}
	for _, route := range router.Routes() {
		fmt.Printf("📦 ROUTE REGISTERED: %s %s\n", route.Method, route.Path)
	}

	router.Run(":8080")
}
