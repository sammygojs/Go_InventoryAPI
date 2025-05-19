package handlers

import (
    "context"
    "log"
    "net/http"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    "github.com/gin-gonic/gin"
	"ProductsAPI/internal/models"
)

func QueryDB(c *gin.Context) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Println("Failed to load AWS config:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to AWS"})
        return
    }

    db := dynamodb.NewFromConfig(cfg)

    out, err := db.Scan(context.TODO(), &dynamodb.ScanInput{
        TableName: awsString("Products"),
    })
    if err != nil {
        log.Println("Scan failed:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query table"})
        return
    }

    var products []*models.Product
    err = attributevalue.UnmarshalListOfMaps(out.Items, &products)
    if err != nil {
        log.Println("Unmarshal failed:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
        return
    }

    c.JSON(http.StatusOK, products)
}

func awsString(s string) *string {
    return &s
}
