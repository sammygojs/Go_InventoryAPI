package handlers

import (
	"testing"
	"ProductsAPI/internal/models"
	"ProductsAPI/internal/services"
)

func TestProductMatchesFilters(t *testing.T) {
	product := &models.Product{
		Variants: []*models.Variant{
			{
				Prices: models.PriceInfo{
					Price: 119.99,
					MembershipPrice: 99.99,
					CurrencyCode: "GBP",
				},
				Inventory: struct {
					Count     interface{} `json:"count"`
					IsInStock bool        `json:"isInStock"`
				}{IsInStock: true},
			},
		},
		Colours: []models.Colour{
			{Colour: "Red"},
			{Colour: "Black"},
		},
	}

	tests := []struct {
		min, max float64
		inStock  bool
		colour   string
		want     bool
	}{
		{90, 110, true, "red", false},
		{101, 200, true, "red", true},
		{90, 110, true, "blue", false},
	}

	for _, tt := range tests {
		got := services.ProductMatchesFilters(product, tt.min, tt.max, tt.inStock, tt.colour)
		if got != tt.want {
			t.Errorf("Failed for input min=%v max=%v stock=%v colour=%v", tt.min, tt.max, tt.inStock, tt.colour)
		}
	}
}
