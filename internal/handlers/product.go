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
	"sort" 
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	startID, _ := strconv.Atoi(c.DefaultQuery("startID", "0")) // optional

	if minPrice == 0 {
		minPrice = 0
	}
	if maxPrice == 0 {
		maxPrice = 999999
	}
	if limit <= 0 {
		limit = 10
	}

	cachedProducts, err := utils.LoadProductsFromDynamo()
	if err != nil {
		log.Printf("Error loading products from Dynamo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load products"})
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

	filtered := make([]*models.Product, 0, len(cachedProducts.Products))
	for _, p := range cachedProducts.Products {
		var clone models.Product
		data, _ := json.Marshal(p)
		_ = json.Unmarshal(data, &clone)

		if locale != "" {
			ApplyTranslation(&clone, locale)
		}
		ApplyMembershipPricing(&clone, isMember)

		if !ProductMatchesFilters(&clone, minPrice, maxPrice, inStock, colourFilter) {
			continue
		}

		filtered = append(filtered, &clone)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ID < filtered[j].ID
	})

	startIndex := 0
	if startID > 0 {
		for i, p := range filtered {
			if p.ID == startID {
				startIndex = i + 1
				break
			}
		}
	}

	endIndex := startIndex + limit
	if endIndex > len(filtered) {
		endIndex = len(filtered)
	}
	page := filtered[startIndex:endIndex]

	var nextStartID *int
	if endIndex < len(filtered) {
		next := filtered[endIndex].ID
		nextStartID = &next
	}

	c.JSON(http.StatusOK, gin.H{
		"count":       len(page),
		"total":       len(filtered),
		"products":    page,
		"hasMore":     nextStartID != nil,
		"nextStartID": nextStartID,
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
		log.Printf("Error fetching product: %v", err)
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
