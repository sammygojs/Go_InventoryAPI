package utils

import (
	"context"
	"fmt"
	"ProductsAPI/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func LoadProductsFromDynamo(minPrice, maxPrice float64, inStock bool, colour string) (*models.Products, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	filterExpr := "price >= :min AND price <= :max"
	exprAttrVals := map[string]types.AttributeValue{
		":min": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", minPrice)},
		":max": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", maxPrice)},
	}

	if inStock {
		filterExpr += " AND inStock = :stock"
		exprAttrVals[":stock"] = &types.AttributeValueMemberBOOL{Value: true}
	}
	if colour != "" {
		filterExpr += " AND contains(colour, :colour)"
		exprAttrVals[":colour"] = &types.AttributeValueMemberS{Value: colour}
	}

	out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String("ProductsTable"),
		FilterExpression:          aws.String(filterExpr),
		ExpressionAttributeValues: exprAttrVals,
	})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to scan DynamoDB: %w", err)
	}

	var productList []*models.Product
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &productList); err != nil {
		return nil, fmt.Errorf("❌ Failed to unmarshal products: %w", err)
	}

	return &models.Products{
		Count:    len(productList),
		Total:    len(productList),
		Products: productList,
	}, nil
}

func LoadSingleProductFromDynamo(id int) (*models.Product, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("ProductsTable"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to get product from DynamoDB: %w", err)
	}
	if out.Item == nil || len(out.Item) == 0 {
		return nil, nil
	}

	var product models.Product
	if err := attributevalue.UnmarshalMap(out.Item, &product); err != nil {
		return nil, fmt.Errorf("❌ Failed to unmarshal product: %w", err)
	}

	return &product, nil
}

// func LoadProductsFromDynamo() (*models.Products, error) {
// 	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
// 	if err != nil {
// 		return nil, fmt.Errorf("❌ Failed to load AWS config: %w", err)
// 	}

// 	client := dynamodb.NewFromConfig(cfg)

// 	out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
// 		TableName: aws.String("ProductsTable"),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("❌ Failed to scan DynamoDB: %w", err)
// 	}

// 	var productList []*models.Product
// 	if err := attributevalue.UnmarshalListOfMaps(out.Items, &productList); err != nil {
// 		return nil, fmt.Errorf("❌ Failed to unmarshal products: %w", err)
// 	}

// 	return &models.Products{
// 		Count:    len(productList),
// 		Total:    len(productList),
// 		Products: productList,
// 	}, nil
// }
