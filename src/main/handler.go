package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type CoffeeShop struct {
	DB     *gorm.DB
	Config Config
}

func (cs *CoffeeShop) handleBuyCoffee(c *gin.Context) {
	userIDHeader, err := strconv.Atoi(c.GetHeader("User-Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid User-Id format"})
		return
	}
	userID := uint(userIDHeader)
	membershipType := c.GetHeader("Membership-Type")
	var coffeeType string
	var coffeeRequest CoffeeRequest
	if c.BindJSON(&coffeeRequest) == nil {
		coffeeType = coffeeRequest.CoffeeType
	}

	quota, found := cs.Config.Memberships[membershipType].Quotas[coffeeType]
	if !found {
		c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid membership type or coffee type"})
		return
	}

	currentTime, waitTime, exceeded := cs.checkQuotaLimits(quota, userID, coffeeType)
	if exceeded {
		c.Header("Retry-After", fmt.Sprintf("%.0f", waitTime.Seconds()))
		c.JSON(http.StatusTooManyRequests, map[string]any{"error": "Quota limit exceeded"})
		return
	}

	usage := QuotaUsage{
		UserID:    userID,
		Coffee:    coffeeType,
		Timestamp: currentTime,
	}
	cs.DB.Create(&usage)

	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, map[string]any{"message": fmt.Sprintf("Enjoy your '%s'", coffeeType)})
}

func (cs *CoffeeShop) checkQuotaLimits(quota Quota, userID uint, coffeeType string) (time.Time, time.Duration, bool) {
	currentTime := time.Now()
	timeThreshold := currentTime.Add(-time.Duration(quota.Retention) * time.Hour)

	var count int
	cs.DB.Model(&QuotaUsage{}).Where("user_id = ? AND coffee = ? AND timestamp > ?", userID, coffeeType, timeThreshold).Count(&count)

	if count >= quota.Amount {
		var waitTime time.Duration
		if lastUsage, err := cs.getLastQuotaUsage(userID, coffeeType); err == nil {
			waitTime = lastUsage.Timestamp.Add(time.Duration(quota.Retention) * time.Hour).Sub(currentTime)
		} else {
			waitTime = 0
		}

		return time.Time{}, waitTime, true
	}
	return currentTime, 0, false
}

func (cs *CoffeeShop) getLastQuotaUsage(userID uint, coffeeType string) (QuotaUsage, error) {
	var usage QuotaUsage
	if err := cs.DB.Model(&QuotaUsage{}).Where("user_id = ? AND coffee = ?", userID, coffeeType).Order("timestamp desc").First(&usage).Error; err != nil {
		return usage, err
	}
	return usage, nil
}
