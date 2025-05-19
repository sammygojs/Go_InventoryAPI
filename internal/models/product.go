package models

import (
	"encoding/json"
	"fmt"
)
type Colour struct {
	SKU    string `json:"sku"`
	Colour string `json:"colour"`
}

type Translation struct {
	ID                 string   `json:"id"`
	Description        string   `json:"description"`
	ShortDescription   string   `json:"shortDescription"`
	Features           []string `json:"features"`
	DefaultCountryCode string   `json:"defaultCountryCode"`
}

type Products struct {
	Count    int        `json:"count"`
	Total    int        `json:"total"`
	Products []*Product `json:"products"`
}

type Product struct {
	ID               int       `dynamodbav:"id" json:"id"`
	SKU              string    `json:"sku"`
	Name             string    `json:"name"`
	ShortDescription *string   `json:"shortDescription"`
	LongDescription  *string   `json:"longDescription"`
	Brand            string    `json:"brand"`
	Price            *Money    `json:"price"`
	Features         []string  `json:"features,omitempty"`
	Images 			 ImageList `json:"images"`
	Options          []*Option `json:"options"`
	Variants         []*Variant  `json:"variants"`
	Translations     []Translation `json:"translations"`
	Colours          []Colour      `json:"colours"`
}

type Variant struct {
	ID    int    `json:"id"`
	EAN   string `json:"ean"`
	SKU   string `json:"sku"`
	Prices struct {
		Price           float64     `json:"price"`
		MembershipPrice interface{} `json:"membershipPrice"`
		CurrencyCode    string      `json:"currencyCode"`
	} `json:"prices"`
	Inventory struct {
		Count     interface{} `json:"count"`
		IsInStock bool        `json:"isInStock"`
	} `json:"inventory"`
	Options []struct {
		ID    int `json:"id"`
		Value []struct {
			Label        string `json:"label"`
			PreorderDate string `json:"preorderDate"`
		} `json:"value"`
		Group string `json:"group"`
	} `json:"options"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Image struct {
	Url     string  `json:"url"`
	AltText *string `json:"altText"`
}

type Option struct {
	ID      string `json:"id"`
	SKU     string `json:"sku"`
	Price   *Money `json:"price"`
	InStock bool   `json:"inStock"`
	Size    string `json:"size"`
	Colour  string `json:"colour"`
}

type ImageList []*Image

func (il *ImageList) UnmarshalJSON(data []byte) error {
	// Attempt 1: Parse as []string
	var urls []string
	if err := json.Unmarshal(data, &urls); err == nil {
		var list []*Image
		for _, u := range urls {
			list = append(list, &Image{Url: u, AltText: nil})
		}
		*il = list
		return nil
	}

	// Attempt 2: Parse as []Image
	var objects []*Image
	if err := json.Unmarshal(data, &objects); err == nil {
		*il = objects
		return nil
	}

	// Fallback: return error
	return fmt.Errorf("ImageList format is invalid")
}