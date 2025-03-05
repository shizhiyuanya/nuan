package vm

import "github.com/jinzhu/gorm"

type Word struct {
	gorm.Model
	English string `json:"english"`
	Chinese string `json:"chinese"`
}

type PaginationParams struct {
	gorm.Model
	Page     string `json:"page"`
	PageSize string `json:"pageSize"`
}
