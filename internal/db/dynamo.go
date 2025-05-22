package db

import (
	"context"
	"fmt"
	"ProductsAPI/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"os"
)

func LoadProductsFromDynamo() (*models.Products, error) {
	if os.Getenv("USE_MOCK_PRODUCTS") == "true" {
		return mockProducts(), nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("Failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		// TableName: aws.String("ProductsTable"),
		TableName: aws.String("TestProductsDB"),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to scan DynamoDB: %w", err)
	}

	// I expect DynamoDB to return many products, and I want each of them stored as a pointer in this list.
	var productList []*models.Product
	// Take this messy JSON-ish DynamoDB output and decode it into []*models.Product
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &productList); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal products: %w", err)
	}

	// jsonData, _ := json.MarshalIndent(productList, "", "  ")
	// log.Println("[DEBUG] Full Product List:\n", string(jsonData))

	return &models.Products{
		Count:    len(productList),
		Total:    len(productList),
		Products: productList,
	}, nil
}


func LoadSingleProductFromDynamo(id int) (*models.Product, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("Failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		// TableName: aws.String("ProductsTable"),
		TableName: aws.String("TestProductsDB"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to get product from DynamoDB: %w", err)
	}
	if out.Item == nil || len(out.Item) == 0 {
		return nil, nil
	}

	var product models.Product
	if err := attributevalue.UnmarshalMap(out.Item, &product); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal product: %w", err)
	}

	return &product, nil
}

func mockProducts() *models.Products {
	return &models.Products{
		Count: 1,
		Total: 1,
		Products: []*models.Product{
			{
				ID:    1,
				SKU:   "MOCK123",
				ShortDescription: ptr("Mock Description"),
				Colours: []models.Colour{
					{SKU: "MOCK123", Colour: "Red/Black"},
				},
				Variants: []*models.Variant{
					{
						ID: 1,
						SKU: "MOCK123",
						Prices: struct {
							Price           float64     `json:"price"`
							MembershipPrice interface{} `json:"membershipPrice"`
							CurrencyCode    string      `json:"currencyCode"`
						}{Price: 119.99, MembershipPrice: 99.99, CurrencyCode: "GBP"},
						Inventory: struct {
							Count     interface{} `json:"count"`
							IsInStock bool        `json:"isInStock"`
						}{IsInStock: true},
					},
				},
				Translations: []models.Translation{
					{
						DefaultCountryCode: "en-gb",
						Description:        "Mock translated description",
						ShortDescription:   "Mock short",
						Features:           []string{"Feature A"},
					},
				},
			},
		},
	}
}

func ptr(s string) *string {
	return &s
}
