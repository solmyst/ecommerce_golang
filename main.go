package main

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.CartItem{}, &models.Order{},
		&models.OrderItem{})

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run()
}
