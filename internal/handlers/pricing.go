package handlers

import (
	"ProductsAPI/internal/models"
	"encoding/json"
	"strconv"
)

func ApplyMembershipPricing(p *models.Product, isMember bool) {
	if !isMember {
		return
	}

	for _, v := range p.Variants {
		var mPriceFloat float64

		switch val := v.Prices.MembershipPrice.(type) {
		case float64:
			mPriceFloat = val
		case string:
			if parsed, err := strconv.ParseFloat(val, 64); err == nil {
				mPriceFloat = parsed
			}
		case json.Number:
			if parsed, err := val.Float64(); err == nil {
				mPriceFloat = parsed
			}
		default:
			continue
		}

		if mPriceFloat > 0 && mPriceFloat < v.Prices.Price {
			v.Prices.Price = mPriceFloat
		}
	}
}
