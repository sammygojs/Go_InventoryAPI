package handlers

import (
	"testing"
	"ProductsAPI/internal/models"
	"ProductsAPI/internal/handlers"
)

func TestApplyTranslation(t *testing.T) {
	p := &models.Product{
		Translations: []models.Translation{
			{
				DefaultCountryCode: "en-gb",
				ShortDescription:   "English short",
				Description:        "English long",
				Features:           []string{"Feature 1"},
			},
		},
	}

	handlers.ApplyTranslation(p, "en-GB")

	if p.ShortDescription == nil || *p.ShortDescription != "English short" {
		t.Error("Translation not applied correctly to ShortDescription")
	}
	if p.LongDescription == nil || *p.LongDescription != "English long" {
		t.Error("Translation not applied correctly to LongDescription")
	}
}
