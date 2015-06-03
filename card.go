package main

import (
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type Card struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name" valid:"alphanum,required"`
	Code      string     `json:"code" valid:"numeric,required"`
	Pin       string     `json:"pin" valid:"numeric,required"`
	IsActive  bool       `json:"is_active" valid:"required"`
	Schedule  []Schedule `json:"schedule" gorm: "many2many:card_schedules;"`
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

func (h *DBHandler) cardsIndexHandler(rw http.ResponseWriter, req *http.Request) {
	page := getPage(req) - 1
	perPage := getPerPage(req)
	offset := perPage * page

	var cards []Card

	h.db.Limit(perPage).Offset(offset).Find(&cards)

	if cards == nil {
		h.r.JSON(rw, http.StatusOK, make([]int64, 0))
	} else {
		h.r.JSON(rw, http.StatusOK, &cards)
	}
}

func (h *DBHandler) cardshowHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	card := Card{}
	h.db.First(&card, id)
	h.r.JSON(rw, http.StatusOK, &card)
}

func (h *DBHandler) cardCreateHandler(rw http.ResponseWriter, req *http.Request) {
	h.cardsEdit(rw, req, 0)
}

func (h *DBHandler) cardUpdateHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	h.cardsEdit(rw, req, id)
}

func (h *DBHandler) cardDeleteHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	card := Card{}
	h.db.Delete(&card, id)
	h.r.JSON(rw, http.StatusOK, &card)
}

func (h *DBHandler) cardsEdit(rw http.ResponseWriter, req *http.Request, id int64) {
	cardForm := CardForm{}

	if err := binding.Bind(req, &cardForm); err.Handle(rw) {
		return
	}

	card := Card{
		Id:       id,
		Name:     cardForm.Name,
		Code:     cardForm.Code,
		Pin:      cardForm.Pin,
		IsActive: cardForm.IsActive,
		Schedule: cardForm.Schedule,
	}

	h.db.Save(&card)
	h.r.JSON(rw, http.StatusOK, &card)
}
