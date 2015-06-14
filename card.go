package main

import (
	_ "fmt"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type Card struct {
	ID        int64      `json:"Id"`
	Name      string     `json:"Name" valid:"alphanum,required"`
	Code      string     `json:"Code" valid:"alphanum,required"`
	Pin       string     `json:"Pin" valid:"numeric,required"`
	IsActive  bool       `json:"IsActive" valid:"required"`
	Schedules []Schedule `json:"Schedules" gorm:"many2many:card_schedule;"`
	CreatedAt time.Time  `json:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt"`
	DeletedAt time.Time  `json:"DeletedAt"`
}

type CardForm struct {
	Name      string
	Code      string
	Pin       string
	IsActive  bool
	Schedules []Schedule
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
		&cf.Schedules: binding.Field{
			Form:     "schedules",
			Required: true,
		},
	}
}

func (h *DBHandler) cardsIndexHandler(rw http.ResponseWriter, req *http.Request) {
	/*
		page := getPage(req) - 1
		perPage := getPerPage(req)
		offset := perPage * page
	*/

	var cards []Card

	h.db.Find(&cards)

	if cards == nil {
		h.r.JSON(rw, http.StatusOK, make([]int64, 0))
	} else {
		h.r.JSON(rw, http.StatusOK, &cards)
	}
}

func (h *DBHandler) cardShowHandler(rw http.ResponseWriter, req *http.Request) {
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

	// lookup the schedule to see if we have it
	// then populate it form our data to avoid an update
	scheduleIds := cardForm.Schedules

	hydratedSchedules := make([]Schedule, len(scheduleIds))
	for _, val := range scheduleIds {
		schedule := Schedule{}
		h.db.First(&schedule, val.ID)

		hydratedSchedules = append(hydratedSchedules, schedule)
	}

	card := Card{
		ID:        id,
		Name:      cardForm.Name,
		Code:      cardForm.Code,
		Pin:       cardForm.Pin,
		IsActive:  cardForm.IsActive,
		Schedules: hydratedSchedules,
	}

	h.db.Save(&card)
	h.r.JSON(rw, http.StatusOK, &card)
}
