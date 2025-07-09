package Models

import (
	"time"
	"gorm.io/gorm"
)

type SiteStatus string

const(
	StatusUp SiteStatus="up"
	StatusDown SiteStatus="down"
	StatusUnknown SiteStatus="unknown"
)

type Sites struct{
	gorm.Model
	//For reminders: `gorm:not null` makes sure the fields are never empty
	UserID uint `gorm:"not null"`
	User User   `gorm:"foreignKey:UserID"`
	Name string `gorm:"not null"`
    URL string `gorm:"not null"`
	CheckInterval int `gorm:"default:10" `
	LastChecked *time.Time
	LastStatus SiteStatus  `gorm:"default:unknown" `
	IsActive bool `gorm:"column:is_active;default:true"`
}

type SiteCheck struct{
	gorm.Model
	SiteID int `gorm:"not null"`
	SiteStatus SiteStatus `gorm:"not null"`
	StatusCode int
	Latency int64
	Error string
}