package dto

import "github.com/jinzhu/gorm"

// 测试接口
type User struct {
	gorm.Model
	Username  string `gorm:"varchar(20);not null"`
	Telephone string `gorm:"varchar(20);not null;unique"`
	Password  string `gorm:"varchar(20);not null;"`
	Avatar    string `gorm:"varchar(255);not null;default:'https://img1.baidu.com/it/u=3736366797,3444523544&fm=253&fmt=auto&app=138&f=JPEG?w=800&h=806'"`
}
