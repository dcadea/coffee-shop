package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	config, err := loadConfig("./src/static/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=5432 user=postgres dbname=coffee_shop sslmode=disable", getDBHost())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	db.AutoMigrate(&QuotaUsage{})

	coffeeShop := &CoffeeShop{
		DB:     db,
		Config: config,
	}

	router := gin.Default()
	router.POST("/coffee", coffeeShop.handleBuyCoffee)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
