package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ProductsAPI/internal/router"
)

// Helper function to perform request
func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Test getting all products
func TestGetAllProducts(t *testing.T) {
	r := router.SetupRouter()
	w := performRequest(r, "GET", "/api/products")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "products")
}

// Test filtering by minimum price
func TestFilterProductsByMinPrice(t *testing.T) {
	r := router.SetupRouter()

	req, _ := http.NewRequest("GET", "/api/products?minPrice=50", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Products []map[string]interface{} `json:"products"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Greater(t, len(resp.Products), 0, "Expected at least one product with price >= 50")

	for _, p := range resp.Products {
		price := p["variants"].([]interface{})[0].(map[string]interface{})["prices"].(map[string]interface{})["price"].(float64)
		assert.GreaterOrEqual(t, price, 50.0, "Product price should be >= 50")
	}
}


// Structs to decode colour field
type Colour struct {
	Colour string `json:"colour"`
}

type Product struct {
	Colours []Colour `json:"colours"`
}

type Response struct {
	Products []Product `json:"products"`
}

func TestFilterProductsByColour(t *testing.T) {
	r := router.SetupRouter()
	w := performRequest(r, "GET", "/api/products?colour=red")

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	for _, product := range resp.Products {
		found := false
		for _, colour := range product.Colours {
			parts := strings.Split(strings.ToLower(colour.Colour), "/")
			for _, part := range parts {
				if strings.TrimSpace(part) == "red" {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		assert.True(t, found, "Expected colour 'red' in product colours")
	}
}

// Test invalid minPrice input
func TestInvalidMinPrice(t *testing.T) {
	r := router.SetupRouter()
	w := performRequest(r, "GET", "/api/products?minPrice=-10")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid")
}
