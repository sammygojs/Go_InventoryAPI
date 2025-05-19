package handlers

import (
	"strings"
	"ProductsAPI/internal/models"
)

func ApplyTranslation(p *models.Product, locale string) {
	// Normalize locale: "de-DE" → "de-de"
	locale = strings.ToLower(locale)

	for _, t := range p.Translations {
		if strings.ToLower(t.DefaultCountryCode) == locale {
			// Override fields
			p.ShortDescription = &t.ShortDescription
			p.LongDescription = &t.Description
			p.Features = t.Features
			break
		}
	}
}
