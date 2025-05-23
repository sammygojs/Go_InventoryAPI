package handlers

import (
	"testing"

	"ProductsAPI/internal/models"
	"ProductsAPI/internal/services"
)

func TestApplyTranslation(t *testing.T) {
	original := &models.Product{
		Translations: []models.Translation{
			{
				DefaultCountryCode: "en-gb",
				ShortDescription:   "English short",
				Description:        "English long",
				Features:           []string{"Feature 1", "Feature 2"},
			},
			{
				DefaultCountryCode: "de-de",
				ShortDescription:   "German short",
				Description:        "German long",
				Features:           []string{"Funktion 1"},
			},
		},
	}

	tests := []struct {
		name         string
		locale       string
		expectedDesc string
		expectedLong string
		expectedFeat []string
	}{
		{
			name:         "Exact match lowercase",
			locale:       "en-gb",
			expectedDesc: "English short",
			expectedLong: "English long",
			expectedFeat: []string{"Feature 1", "Feature 2"},
		},
		{
			name:         "Exact match uppercase",
			locale:       "EN-GB",
			expectedDesc: "English short",
			expectedLong: "English long",
			expectedFeat: []string{"Feature 1", "Feature 2"},
		},
		{
			name:         "Different match",
			locale:       "de-DE",
			expectedDesc: "German short",
			expectedLong: "German long",
			expectedFeat: []string{"Funktion 1"},
		},
		{
			name:         "No match found",
			locale:       "fr-FR",
			expectedDesc: "", // Should remain nil
			expectedLong: "", // Should remain nil
			expectedFeat: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &models.Product{
				Translations: original.Translations,
			}

			services.ApplyTranslation(p, tt.locale)

			if tt.expectedDesc != "" {
				if p.ShortDescription == nil || *p.ShortDescription != tt.expectedDesc {
					t.Errorf("ShortDescription mismatch: got %v, want %v", p.ShortDescription, tt.expectedDesc)
				}
			} else if p.ShortDescription != nil {
				t.Errorf("Expected nil ShortDescription, got %v", *p.ShortDescription)
			}

			if tt.expectedLong != "" {
				if p.LongDescription == nil || *p.LongDescription != tt.expectedLong {
					t.Errorf("LongDescription mismatch: got %v, want %v", p.LongDescription, tt.expectedLong)
				}
			} else if p.LongDescription != nil {
				t.Errorf("Expected nil LongDescription, got %v", *p.LongDescription)
			}

			if len(tt.expectedFeat) > 0 {
				if len(p.Features) != len(tt.expectedFeat) {
					t.Errorf("Features length mismatch: got %v, want %v", len(p.Features), len(tt.expectedFeat))
				}
				for i := range p.Features {
					if p.Features[i] != tt.expectedFeat[i] {
						t.Errorf("Feature[%d] mismatch: got %v, want %v", i, p.Features[i], tt.expectedFeat[i])
					}
				}
			} else if p.Features != nil {
				t.Errorf("Expected nil Features, got %v", p.Features)
			}
		})
	}
}
