package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	email    string `gorm:"uniqueIndex"`
	Password string
}

type Site struct {
	gorm.Model
	Url       string `form:"url"`
	Endpoints []Endpoint
}

type Endpoint struct {
	gorm.Model
	Path        string        `form:"path"`
	Status      string        `form:"status"`
	Frequency   time.Duration `form:"frequency"`
	LastChecked time.Time     `form:"last_checked"`
	Uptime      float32       `form:"uptime"`
	SiteID      uint          `form:"site_id"`
}
