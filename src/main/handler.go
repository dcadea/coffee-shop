package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	if err := cs.lockUserQuotasUsage(userID, coffeeType); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to process the request"})
		return
	}

	usage := QuotaUsage{
		UserID:    userID,
		Coffee:    coffeeType,
		Timestamp: currentTime,
	}
	cs.DB.Create(&usage)
	cs.unlockUserQuotasUsage(userID, coffeeType)

	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, map[string]any{"message": fmt.Sprintf("Enjoy your '%s'", coffeeType)})
}

func (cs *CoffeeShop) checkQuotaLimits(quota Quota, userID uint, coffeeType string) (time.Time, time.Duration, bool) {
	currentTime := time.Now()
	timeThreshold := currentTime.Add(-time.Duration(quota.Retention) * time.Hour)

	var count int64
	cs.DB.Model(&QuotaUsage{}).Where("user_id = ? AND coffee = ? AND timestamp > ?", userID, coffeeType, timeThreshold).Count(&count)

	if count >= int64(quota.Amount) {
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

func (cs *CoffeeShop) lockUserQuotasUsage(userID uint, coffeeType string) error {
	query := fmt.Sprintf("SELECT pg_advisory_lock(hashtext('%d'), hashtext('%s'))", userID, coffeeType)
	result := cs.DB.Exec(query)
	if result.Error != nil || result.RowsAffected != 1 {
		return fmt.Errorf("failed to acquire lock")
	}
	return nil
}

func (cs *CoffeeShop) unlockUserQuotasUsage(userID uint, coffeeType string) {
	query := fmt.Sprintf("SELECT pg_advisory_unlock(hashtext('%d'), hashtext('%s'))", userID, coffeeType)
	_ = cs.DB.Exec(query)
}
