package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"ProductsAPI/internal/models"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	// Read products.json
	data, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatalf("❌ Failed to read products.json: %v", err)
	}

	var products []*models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		log.Fatalf("❌ Failed to unmarshal JSON: %v", err)
	}

	// Upload each product to DynamoDB
	for _, product := range products {
		log.Printf("Uploading product: ID=%d, SKU=%s", product.ID, product.SKU)
		item, err := attributevalue.MarshalMap(product)
		if err != nil {
			log.Fatalf("❌ Failed to marshal product: %v", err)
		}

		_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("ProductsTable"),
			Item:      item,
		})
		if err != nil {
			log.Fatalf("❌ Failed to insert product ID %d: %v", product.ID, err)
		}

		log.Printf("✅ Uploaded product ID: %d", product.ID)
	}
}