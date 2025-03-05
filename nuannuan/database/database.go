package database

import (
	"fmt"
	"nuannuan/model/dto"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "数据库IP"
	port := "7788"
	database := "数据库表名"
	username := "数据库用户名"
	password := "数据库用户密码"
	charset := "utf8mb4"
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}

	// 自动生成技术
	db.AutoMigrate(&dto.User{})
	db.AutoMigrate(&dto.Picture{})
	db.AutoMigrate(&dto.Tag{})
	db.AutoMigrate(&dto.Word{})
	db.AutoMigrate(&dto.Video{})
	db.AutoMigrate(&dto.VideoType{})
	db.AutoMigrate(&dto.Message{})
	DB = db

	return db

}

func GetDB() *gorm.DB {
	return DB
}
