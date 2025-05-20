package handlers

import (
	"strings"
	"ProductsAPI/internal/models"
)

func ProductMatchesFilters(p *models.Product, minPrice, maxPrice float64, requireInStock bool, colour string) bool {
	matchesPrice := false
	matchesStock := false

	for _, v := range p.Variants {
		price := v.Prices.Price

		if (minPrice == 0 || price >= minPrice) &&
			(maxPrice == 0 || price <= maxPrice) {
			matchesPrice = true
		}

		if !requireInStock || v.Inventory.IsInStock {
			matchesStock = true
		}
	}

	if !matchesPrice || !matchesStock {
		return false
	}

	if colour != "" {
		found := false
		for _, c := range p.Colours {
			parts := strings.Split(strings.ToLower(c.Colour), "/")
			for _, part := range parts {
				if part == colour {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}
	

	return true
}
