package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PlaceOrder(c *gin.Context) {
	userRaw, _ := c.Get("user")
	user := userRaw.(models.User)

	var cartItems []models.CartItem
	if err := database.DB.
		Where("user_id = ?", user.ID).
		Preload("Product").
		Find(&cartItems).Error; err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	tx := database.DB.Begin()

	order := models.Order{UserID: user.ID}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create order"})
		return
	}

	for _, item := range cartItems {
		if item.Product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock for " + item.Product.Name,
			})
			return
		}

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order item"})
			return
		}

		// Deduct stock
		item.Product.Stock -= item.Quantity
		tx.Save(&item.Product)
	}

	// Clear cart
	if err := tx.Where("user_id = ?", user.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order_id": order.ID})
}

func GetOrderByID(c *gin.Context) {
	userRaw, _ := c.Get("user")
	user := userRaw.(models.User)
	orderID, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := database.DB.Preload("Items.Product").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.UserID != user.ID && !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetOrders(c *gin.Context) {
	userRaw, _ := c.Get("user")
	user := userRaw.(models.User)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var orders []models.Order
	query := database.DB.Preload("Items.Product").Limit(limit).Offset(offset)

	if !user.IsAdmin {
		query = query.Where("user_id = ?", user.ID)
	}

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
