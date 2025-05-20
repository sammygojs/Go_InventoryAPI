package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"ProductsAPI/internal/handlers"
	"ProductsAPI/internal/models"
)

func TestGetProductsIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/products", handlers.GetProducts)

	// Mock data
	handlers.SetCachedProducts(mockProducts()) 

	// Create request
	req, _ := http.NewRequest("GET", "/api/products?limit=2&colour=red", nil)
	req.Header.Set("Accept-Language", "en-gb")
	req.Header.Set("X-Member", "true")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 but got %d", w.Code)
	}

	if !contains(w.Body.String(), "products") {
		t.Errorf("Expected 'products' in response body, got: %s", w.Body.String())
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func mockProducts() *models.Products {
	return &models.Products{
		Count: 1,
		Total: 1,
		Products: []*models.Product{
			{
				ID:    1,
				SKU:   "MOCK123",
				Brand: "MockBrand",
				ShortDescription: ptr("Mock product"),
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
						Description:        "Translated Description",
						ShortDescription:   "Translated Short",
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
