package main

import (
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
)

type DBHandler struct {
	db *gorm.DB
	r  *render.Render
}

type PaginatedResponse struct {
	Items      []interface{} `json:"items"`
	TotalItems int           `json:"total_items"`
}
