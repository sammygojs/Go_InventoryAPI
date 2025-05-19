package tests

import (
	"ProductsAPI/internal/models"
	"ProductsAPI/internal/handlers"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyMembershipPricing(t *testing.T) {
	p := &models.Product{
		Variants: []*models.Variant{
			{
				ID: 1,
				Prices: struct {
					Price           float64     `json:"price"`
					MembershipPrice interface{} `json:"membershipPrice"`
					CurrencyCode    string      `json:"currencyCode"`
				}{
					Price:           100.0,
					MembershipPrice: json.Number("79.99"),
					CurrencyCode:    "USD",
				},
			},
		},
	}

	handlers.ApplyMembershipPricing(p, true)

	assert.Equal(t, 79.99, p.Variants[0].Prices.Price)
}

func TestApplyMembershipPricing_NoChangeIfHigher(t *testing.T) {
	p := &models.Product{
		Variants: []*models.Variant{
			{
				ID: 2,
				Prices: struct {
					Price           float64     `json:"price"`
					MembershipPrice interface{} `json:"membershipPrice"`
					CurrencyCode    string      `json:"currencyCode"`
				}{
					Price:           100.0,
					MembershipPrice: 150.0, // higher than base price
					CurrencyCode:    "USD",
				},
			},
		},
	}

	handlers.ApplyMembershipPricing(p, true)

	assert.Equal(t, 100.0, p.Variants[0].Prices.Price)
}
