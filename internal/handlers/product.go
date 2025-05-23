package handlers

import (
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

var cachedProducts *models.Products

func SetCachedProducts(p *models.Products) {
	cachedProducts = p
}

// GET /api/products
func GetProducts(c *gin.Context) {
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")

	minPrice := 0.0
	if minPriceStr != "" {
		var err error
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil || minPrice < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid minPrice"})
			return
		}
	}

	maxPrice := 999999.0
	if maxPriceStr != "" {
		var err error
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil || maxPrice < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maxPrice"})
			return
		}
	}

	inStockQuery := c.Query("inStock")
	requireInStock := inStockQuery == "true"
	stockFilterApplied := inStockQuery != ""

	colourFilter := strings.ToLower(c.Query("colour"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	productsData, err := db.LoadProductsFromDynamo()
	if err != nil {
		log.Printf("[ERROR] Failed to load products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Problem loading products"})
		return
	}
	cachedProducts = productsData

	locale := c.GetHeader("Accept-Language")
	if locale == "" {
		locale = c.Query("locale")
	}
	if len(locale) > 5 {
		locale = locale[:5]
	}
	isMember := c.GetHeader("X-Member") == "true"

	filtered := make([]*models.Product, 0, len(cachedProducts.Products))

	for _, p := range cachedProducts.Products {
		if locale != "" {
			services.ApplyTranslation(p, locale)
		}
		services.ApplyMembershipPricing(p, isMember)

		if !services.ProductMatchesFilters(p, minPrice, maxPrice, requireInStock, stockFilterApplied, colourFilter) {
			continue
		}
		filtered = append(filtered, p)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ID < filtered[j].ID
	})

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

	c.JSON(http.StatusOK, gin.H{
		"count":    len(pageSlice),
		"total":    len(filtered),
		"page":     page,
		"nextPage": nextPage,
		"hasMore":  hasMore,
		"products": pageSlice,
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

	product, err := db.LoadSingleProductFromDynamo(id)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't load product"})
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

	if locale != "" {
		services.ApplyTranslation(product, locale)
	}
	services.ApplyMembershipPricing(product, isMember)

	c.JSON(http.StatusOK, product)
}
