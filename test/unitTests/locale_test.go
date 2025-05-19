package tests

import (
	"ProductsAPI/internal/models"
	"ProductsAPI/internal/handlers"
	"testing"
)

func TestApplyTranslation(t *testing.T) {
	product := &models.Product{
		ShortDescription: strPtr("Default Short Description"),
		LongDescription:  strPtr("Default Long Description"),
		Features:         []string{"Default Feature"},
		Translations: []models.Translation{
			{
				ID:                 "GERMAN",
				DefaultCountryCode: "de-de",
				ShortDescription:   "Kurzbeschreibung",
				Description:        "Lange Beschreibung",
				Features:           []string{"Merkmal 1", "Merkmal 2"},
			},
			{
				ID:                 "ENGLISH",
				DefaultCountryCode: "en-gb",
				ShortDescription:   "Short Desc",
				Description:        "Long Desc",
				Features:           []string{"Feature A", "Feature B"},
			},
		},
	}

	tests := []struct {
		locale         string
		expectedShort  string
		expectedLong   string
		expectedFirstFeature string
	}{
		{"de-DE", "Kurzbeschreibung", "Lange Beschreibung", "Merkmal 1"},
		{"en-GB", "Short Desc", "Long Desc", "Feature A"},
		{"fr-FR", "Default Short Description", "Default Long Description", "Default Feature"},
	}

	for _, tc := range tests {
		// Reset product fields to defaults before each test
		product.ShortDescription = strPtr("Default Short Description")
		product.LongDescription = strPtr("Default Long Description")
		product.Features = []string{"Default Feature"}

		handlers.ApplyTranslation(product, tc.locale)

		if product.ShortDescription == nil || *product.ShortDescription != tc.expectedShort {
			t.Errorf("[%s] Expected shortDescription %q, got %q", tc.locale, tc.expectedShort, safeStr(product.ShortDescription))
		}

		if product.LongDescription == nil || *product.LongDescription != tc.expectedLong {
			t.Errorf("[%s] Expected longDescription %q, got %q", tc.locale, tc.expectedLong, safeStr(product.LongDescription))
		}

		if len(product.Features) == 0 || product.Features[0] != tc.expectedFirstFeature {
			t.Errorf("[%s] Expected first feature %q, got %q", tc.locale, tc.expectedFirstFeature, product.Features)
		}
	}
}

func strPtr(s string) *string {
	return &s
}

func safeStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
