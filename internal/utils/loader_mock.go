//go:build test
// +build test

package utils

import "ProductsAPI/internal/models"

func LoadProductsFromDynamo() (*models.Products, error) {
	return &models.Products{
		Count: 1,
		Total: 1,
		Products: []*models.Product{
			{
				ID:    1,
				SKU:   "MOCK123",
				Brand: "MockBrand",
				ShortDescription: ptr("Mock Product Description"),
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
						}{
							Price: 119.99, MembershipPrice: 99.99, CurrencyCode: "GBP",
						},
						Inventory: struct {
							Count     interface{} `json:"count"`
							IsInStock bool        `json:"isInStock"`
						}{
							Count: nil, IsInStock: true,
						},
						Options: []struct {
							ID    int `json:"id"`
							Value []struct {
								Label        string `json:"label"`
								PreorderDate string `json:"preorderDate"`
							} `json:"value"`
							Group string `json:"group"`
						}{
							{
								ID:    1,
								Group: "Size",
								Value: []struct {
									Label        string `json:"label"`
									PreorderDate string `json:"preorderDate"`
								}{
									{Label: "8", PreorderDate: ""},
								},
							},
						},
					},
				},
				Translations: []models.Translation{
					{
						DefaultCountryCode: "en-gb",
						Description:        "Translated long desc",
						ShortDescription:   "Translated short desc",
						Features:           []string{"Feature A", "Feature B"},
					},
				},
			},
		},
	}, nil
}

func ptr(s string) *string {
	return &s
}
