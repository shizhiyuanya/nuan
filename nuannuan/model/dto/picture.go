package dto

import "github.com/jinzhu/gorm"

type Picture struct {
	gorm.Model
	Name  string `gorm:"varchar(255);not null;unique"`
	TagId string `gorm:"varchar(255);not null"`
}

type Tag struct {
	gorm.Model
	TagName string `gorm:"varchar(20);not null;unique"`
}
