package dto

import (
	"time"
)

type Message struct {
	ID         uint   `gorm:"primary_key"`
	FromUserID string `gorm:"type:varchar(255);not null"`
	ToUserID   string `gorm:"type:varchar(255);not null"`
	Message    string `gorm:"type:text;not null"`
	IsRead     bool   `gorm:"default:false"`
	CreatedAt  time.Time
}
