package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"ProductsAPI/internal/models"
	"ProductsAPI/internal/db"
	"ProductsAPI/internal/services"
)

// Replaces cache every time for simplicity (you could optimize later)
var cachedProducts *models.Products

// Helper to manually set cache if needed
func SetCachedProducts(p *models.Products) {
	cachedProducts = p
}

// GET /api/products
func GetProducts(c *gin.Context) {
	// Grab filters from the query string
	minPrice, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)
	inStock := c.Query("inStock") == "true"
	colourFilter := strings.ToLower(c.Query("colour"))

	// Handle pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Fallbacks for price filters
	if minPrice == 0 {
		minPrice = 0
	}
	if maxPrice == 0 {
		maxPrice = 999999
	}

	// log.Printf("[INFO] Getting products (page: %d, limit: %d, minPrice: %.2f, maxPrice: %.2f, inStock: %v, colour: %s)",
	// 	page, limit, minPrice, maxPrice, inStock, colourFilter)

	// Always pull fresh data from DynamoDB 
	// productsData is of Products Struct
	productsData, err := db.LoadProductsFromDynamo()
	// log.Println("[DEBUG] productsData: ",productsData)

	// serializes your Go struct into a JSON string.
	// data, _ := json.MarshalIndent(productsData, "", "  ")
	// log.Println("[DEBUG] Full Products JSON:\n", string(data))

	if err != nil {
		log.Printf("[ERROR] Failed to load products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem loading products"})
		return
	}
	//cachedProducts to accessed throughout
	cachedProducts = productsData

	// Figure out user's language and membership status
	locale := c.GetHeader("Accept-Language")
	if locale == "" {
		locale = c.Query("locale")
	}
	if len(locale) > 5 {
		locale = locale[:5]
	}
	isMember := c.GetHeader("X-Member") == "true"
	// log.Printf("[DEBUG] Locale: %s | Member: %v", locale, isMember)

	// Start filtering and transforming products
	filtered := make([]*models.Product, 0, len(cachedProducts.Products))
	// log.Println("[DEBUG] filtered:\n", string(filtered))
	// log.Println("[DEBUG] filtered: ",filtered)
	for _, p := range cachedProducts.Products {
		// Deep clone to avoid mutating original
		var clone models.Product

		data, _ := json.Marshal(p)
		// log.Println("data: ",data)
		//returns an error which is ignored
		_ = json.Unmarshal(data, &clone)

		// Apply translations and member pricing
		if locale != "" {
			services.ApplyTranslation(&clone, locale)
		}
		services.ApplyMembershipPricing(&clone, isMember)

		// Only keep products that match the filters
		if !services.ProductMatchesFilters(&clone, minPrice, maxPrice, inStock, colourFilter) {
			continue
		}
		filtered = append(filtered, &clone)
	}

	// log.Printf("[INFO] %d products matched filters", len(filtered))

	// Sort them by ID so pagination works predictably
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ID < filtered[j].ID
	})

	// Pagination logic: calculate offset, boundaries
	offset := (page - 1) * limit
	startIndex := offset
	if startIndex > len(filtered) {
		startIndex = len(filtered)
	}
	endIndex := startIndex + limit
	if endIndex > len(filtered) {
		endIndex = len(filtered)
	}

	pageSlice := filtered[startIndex:endIndex]
	hasMore := endIndex < len(filtered)
	nextPage := page + 1
	if !hasMore {
		nextPage = 0
	}

	// log.Printf("[INFO] Sending page %d (%d items). More: %v", page, len(pageSlice), hasMore)

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"count":     len(pageSlice),
		"total":     len(filtered),
		"page":      page,
		"nextPage":  nextPage,
		"hasMore":   hasMore,
		"products":  pageSlice,
	})
}

// GET /api/products/:productID
func GetProduct(c *gin.Context) {
	idParam := c.Param("productID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("[WARN] Invalid product ID: %s", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	log.Printf("[INFO] Looking up product ID: %d", id)

	product, err := db.LoadSingleProductFromDynamo(id)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't load product"})
		return
	}
	if product == nil {
		log.Printf("[INFO] Product %d not found", id)
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

	log.Printf("[DEBUG] Locale for single product: %s | Member: %v", locale, isMember)

	// Clone to modify freely
	var clone models.Product
	data, _ := json.Marshal(product)
	_ = json.Unmarshal(data, &clone)

	if locale != "" {
		services.ApplyTranslation(&clone, locale)
	}
	services.ApplyMembershipPricing(&clone, isMember)

	c.JSON(http.StatusOK, clone)
}
