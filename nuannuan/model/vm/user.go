package vm

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}
