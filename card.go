package main

import (
	"github.com/mholt/binding"
	"time"
)

type Card struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name" valid:"alphanum,required"`
	Code      string     `json:"code" valid:"numeric,required"`
	Pin       string     `json:"pin" valid:"numeric,required"`
	IsActive  bool       `json:"is_active" valid:"required"`
	Scheudle  []Schedule `json:"schedule" gorm: "many2many:card_schedules;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt time.Time  `json:"deleted_at"`
}

func (c *Card) GetName() string {
	return "card"
}

type CardForm struct {
	Name     string
	Code     string
	Pin      string
	IsActive bool
	Schedule []Schedule
}

func (cf *CardForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&cf.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
		&cf.Code: binding.Field{
			Form:     "code",
			Required: true,
		},
		&cf.Pin: binding.Field{
			Form:     "pin",
			Required: true,
		},
		&cf.IsActive: binding.Field{
			Form:     "is_active",
			Required: true,
		},
		&cf.Schedule: binding.Field{
			Form:     "schedules",
			Required: true,
		},
	}
}
