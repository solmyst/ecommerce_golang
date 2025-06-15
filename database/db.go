package database

import (
	"fmt"
	"log"
	"os"

	"ecommerce/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// You can use env variables or hardcode here for dev
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = db
	fmt.Println("Connected to MySQL database!")
}
func Migrate() {
	DB.AutoMigrate(
		&models.User{},
		// &models.Product{},
		// &models.CartItem{},
		// &models.Order{},
		// &models.OrderItem{}, // if you have one
	)
}
