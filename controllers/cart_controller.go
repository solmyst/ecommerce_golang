package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	userRaw, _ := c.Get("user")
	user := userRaw.(models.User)

	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.Stock < input.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	var existingItem models.CartItem
	if err := database.DB.
		Where("user_id = ? AND product_id = ?", user.ID, input.ProductID).
		First(&existingItem).Error; err == nil {
		// Item exists, update quantity
		existingItem.Quantity += input.Quantity
		database.DB.Save(&existingItem)
		c.JSON(http.StatusOK, existingItem)
		return
	}

	// Else create new cart item
	cartItem := models.CartItem{
		UserID:    user.ID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := database.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add to cart"})
		return
	}

	c.JSON(http.StatusOK, cartItem)
}

func ViewCart(c *gin.Context) {
	userRaw, _ := c.Get("user")
	user := userRaw.(models.User)

	var items []models.CartItem
	if err := database.DB.
		Where("user_id = ?", user.ID).
		Preload("Product").
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch cart"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func RemoveFromCart(c *gin.Context) {
	userRaw, _ := c.Get("user")
	user := userRaw.(models.User)

	itemID, _ := strconv.Atoi(c.Param("item_id"))

	var item models.CartItem
	if err := database.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if item.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your cart item"})
		return
	}

	database.DB.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"message": "Item removed"})
}
