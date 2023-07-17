package main

import (
	"time"
)

type CoffeeRequest struct {
	CoffeeType string `json:"coffee_type"`
}

type Quota struct {
	Amount    int `yaml:"amount"`
	Retention int `yaml:"retention"`
}

type Membership struct {
	Quotas map[string]Quota `yaml:"quotas"`
}

type Config struct {
	Memberships map[string]Membership `yaml:"memberships"`
}

type QuotaUsage struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint
	Coffee    string
	Timestamp time.Time
}
