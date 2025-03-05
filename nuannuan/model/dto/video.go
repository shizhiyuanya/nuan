package dto

import "github.com/jinzhu/gorm"

type Video struct {
	gorm.Model
	Name   string `gorm:"varchar(255);not null;unique"`
	TypeId string `gorm:"varchar(255);not null"`
}

type VideoType struct {
	gorm.Model
	TypeName string `gorm:"varchar(20);not null;unique"`
}
