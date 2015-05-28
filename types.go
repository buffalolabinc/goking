package main

import (
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
)

type DBHandler struct {
	db *gorm.DB
	r  *render.Render
}

type Model interface {
	GetName() string
}
