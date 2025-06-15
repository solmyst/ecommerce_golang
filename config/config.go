package config

import "os"

var (
	DB_URL = os.Getenv("DB_URL") // e.g., "postgres://user:password@localhost:5432/ecommerce"
	PORT   = ":8080"
)
