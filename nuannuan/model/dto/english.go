package dto

import "github.com/jinzhu/gorm"

type Word struct {
	gorm.Model
	English string `gorm:"varchar(255);not null;unique"`
	Chinese string `gorm:"varchar(255);not null;unique"`
}
