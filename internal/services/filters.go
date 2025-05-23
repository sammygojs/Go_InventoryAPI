package services

import (
	"strings"
	"ProductsAPI/internal/models"
)

func ProductMatchesFilters(p *models.Product, minPrice, maxPrice float64, requireInStock bool, stockFilterApplied bool, colour string) bool {
	priceMatch := false
	stockMatch := false

	for _, v := range p.Variants {
		price := v.Prices.Price
		if price >= minPrice && price <= maxPrice {
			priceMatch = true
		}

		if !stockFilterApplied || v.Inventory.IsInStock == requireInStock {
			stockMatch = true
		}
	}

	if !priceMatch || !stockMatch {
		return false
	}

	if colour != "" {
		target := strings.ToLower(colour)
		for _, c := range p.Colours {
			parts := strings.Split(strings.ToLower(c.Colour), "/")
			for _, part := range parts {
				if strings.TrimSpace(part) == target {
					return true
				}
			}
		}
		return false
	}

	return true
}
