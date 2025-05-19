package handlers

import (
	"encoding/json" 
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"ProductsAPI/internal/models"
	"ProductsAPI/internal/utils"
	"strings"
	"log"
)

var cachedProducts *models.Products

func SetCachedProducts(p *models.Products) {
	cachedProducts = p
}

func GetProducts(c *gin.Context) {
	minPrice, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)
	inStock := c.Query("inStock") == "true"
	colourFilter := strings.ToLower(c.Query("colour"))

	// ✅ Apply defaults
	if minPrice == 0 {
		minPrice = 0
	}
	if maxPrice == 0 {
		maxPrice = 999999 // or whatever upper limit you want
	}
	
	cachedProducts, err := utils.LoadProductsFromDynamo(minPrice, maxPrice, inStock, colourFilter)

	if err != nil {
		log.Printf("❌ Error loading products from Dynamo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load products"})
		return
	}

	// Locale detection
	locale := c.GetHeader("Accept-Language")
	if locale == "" {
		locale = c.Query("locale")
	}
	if len(locale) > 5 {
		locale = locale[:5]
	}

	// Membership detection
	isMember := c.GetHeader("X-Member") == "true"

	// Parse filter params
	minPrice, _ = strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPrice, _ = strconv.ParseFloat(c.Query("maxPrice"), 64)
	inStock = c.Query("inStock") == "true"
	colourFilter = strings.ToLower(c.Query("colour"))

	filtered := make([]*models.Product, 0, len(cachedProducts.Products))

	for _, p := range cachedProducts.Products {
		// clone := *p
		var clone models.Product
		data, _ := json.Marshal(p)
		_ = json.Unmarshal(data, &clone)

		if locale != "" {
			ApplyTranslation(&clone, locale)
		}
		ApplyMembershipPricing(&clone, isMember)

		// // Apply filtering
		// if !productMatchesFilters(&clone, minPrice, maxPrice, inStock, colourFilter) {
		// 	continue
		// }

		filtered = append(filtered, &clone)
	}

	c.JSON(http.StatusOK, models.Products{
		Count:    len(filtered),
		Total:    len(filtered),
		Products: filtered,
	})
}

func GetProduct(c *gin.Context) {
	idParam := c.Param("productID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := utils.LoadSingleProductFromDynamo(id)
	if err != nil {
		log.Printf("❌ Error fetching product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load product"})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	locale := c.GetHeader("Accept-Language")
	if locale == "" {
		locale = c.Query("locale")
	}
	if len(locale) > 5 {
		locale = locale[:5]
	}

	isMember := c.GetHeader("X-Member") == "true"

	var clone models.Product
	data, _ := json.Marshal(product)
	_ = json.Unmarshal(data, &clone)

	if locale != "" {
		ApplyTranslation(&clone, locale)
	}
	ApplyMembershipPricing(&clone, isMember)

	c.JSON(http.StatusOK, clone)
}



// func GetProduct(c *gin.Context) {
// 	cachedProducts, err := utils.LoadProductsFromDynamo()
// 	if err != nil {
// 		log.Printf("❌ Error loading products from Dynamo: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load products"})
// 		return
// 	}

// 	idParam := c.Param("productID")
// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
// 		return
// 	}

// 	locale := c.GetHeader("Accept-Language")
// 	if locale == "" {
// 		locale = c.Query("locale")
// 	}
// 	if len(locale) > 5 {
// 		locale = locale[:5]
// 	}

// 	isMember := c.GetHeader("X-Member") == "true"

// 	for _, p := range cachedProducts.Products {
// 		if p.ID == id {

// 			var clone models.Product
// 			data, _ := json.Marshal(p)
// 			_ = json.Unmarshal(data, &clone)

// 			if locale != "" {
// 				applyTranslation(&clone, locale)
// 			}
// 			applyMembershipPricing(&clone, isMember)

// 			c.JSON(http.StatusOK, clone)
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
// }

