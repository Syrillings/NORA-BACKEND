package Models

import (
	"time"
	"gorm.io/gorm"
)

type PasswordResetToken struct{
	gorm.Model 
	Token string `gorm:"unique"`
	UserID string `gorm:"unique"`
	ExpiresAt time.Time 
	Used bool `gorm:"default:false"`
}