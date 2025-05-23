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
					Price:           119.99,
					MembershipPrice: 99.99,
					CurrencyCode:    "GBP",
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
		name              string
		min               float64
		max               float64
		inStock           bool
		stockFilterActive bool
		colour            string
		expected          bool
	}{
		{"Price too low", 90, 110, true, true, "red", false},
		{"Valid range, matching colour", 101, 200, true, true, "red", true},
		{"Invalid colour", 90, 110, true, true, "blue", false},
		{"No colour filter", 90, 200, true, true, "", true},
		{"In stock false but product is in stock", 101, 200, false, true, "red", false},
		{"Upper bound edge case", 119.99, 119.99, true, true, "red", true},
		{"Lower bound edge case", 119.99, 200, true, true, "red", true},
		{"Case insensitive colour", 90, 200, true, true, "ReD", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := services.ProductMatchesFilters(product, tt.min, tt.max, tt.inStock, tt.stockFilterActive, tt.colour)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v (min=%.2f, max=%.2f, inStock=%v, stockFilterActive=%v, colour=%q)",
					tt.expected, result, tt.min, tt.max, tt.inStock, tt.stockFilterActive, tt.colour)
			}
		})
	}
}
