package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func main() {
	config, err := loadConfig("./src/static/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=5432 user=postgres dbname=coffee_shop sslmode=disable", getDBHost()))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

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
