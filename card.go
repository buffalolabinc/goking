package main

import (
	"net/http"
	"time"
	"github.com/asaskevich/govalidator"
	"encoding/json"
)

type CardIdentity struct {
	ID	int64	`json:"id" valid:"required"`
}

type Card struct {
	CardIdentity
	Name      string     `json:"name" valid:"printableascii,required"`
	Pin       string     `json:"-" valid:"numeric,required"`
	IsActive  bool       `json:"is_active"`
	Code      string     `json:"code" valid:"alphanum,required"`
	Schedules []Schedule `json:"schedules" gorm:"many2many:card_schedule;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt time.Time  `json:"deleted_at"`
}

type CreateCardRequest struct {
	Name      string     `json:"name"`
	Pin       string     `json:"pin"`
	IsActive  bool       `json:"is_active"`
	Code      string     `json:"code"`
	Schedules []ScheduleIdentity `json:"schedules"`
}

type UpdateCardRequest struct {
	Name      string     `json:"name"`
	Pin       string     `json:"pin"`
	IsActive  bool       `json:"is_active"`
	Schedules []ScheduleIdentity `json:"schedules"`
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
		var count int
		h.db.Table("cards").Count(&count)

		vals := make([]interface{}, len(cards))
		for i, v := range cards {
			vals[i] = v
		}

		resp := getResponse(vals, count)
		h.r.JSON(rw, http.StatusOK, &resp)
	}
}

func (h *DBHandler) cardShowHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	card := Card{}

	// Find the Card
	h.db.First(&card, id)

	// Hydrate our associations
	h.db.Model(&card).Association("Schedules").Find(&card.Schedules)

	h.r.JSON(rw, http.StatusOK, &card)
}

func (h *DBHandler) cardCreateHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	createCard := CreateCardRequest{}

	err := decoder.Decode(&createCard)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	var cardCount int

	h.db.Model(Card{}).Where("code = ? AND is_active = ?", createCard.Code, true).Count(&cardCount)

	if cardCount > 0 {
		h.r.JSON(rw, http.StatusConflict, map[string]string{"error":"More than one active card with that code"})
		return
	}

	card := Card{}

	card.Name = createCard.Name
	card.Code = createCard.Code
	card.Pin = createCard.Pin
	card.IsActive = createCard.IsActive

	_, err = govalidator.ValidateStruct(&card)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	scheduleIds := []int64{}

	for _, val := range createCard.Schedules {
		scheduleIds = append(scheduleIds, val.ID)
	}

	schedules := []Schedule{}

	h.db.Where("id in (?)", scheduleIds).Find(&schedules)

	card.Schedules = schedules

	h.db.Save(&card)
	
	h.r.JSON(rw, http.StatusOK, &card)
}

func (h *DBHandler) cardUpdateHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	updateCard := UpdateCardRequest{}

	err := decoder.Decode(&updateCard)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	id := getId(req)
	card := Card{}

	h.db.First(&card, id)

	card.Name = updateCard.Name
	card.IsActive = updateCard.IsActive

	if len(updateCard.Pin) > 0 {
		card.Pin = updateCard.Pin
	}

	_, err = govalidator.ValidateStruct(&card)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	h.db.Save(&card)

	scheduleIds := []int64{}

	for _, val := range updateCard.Schedules {
		scheduleIds = append(scheduleIds, val.ID)
	}

	schedules := []Schedule{}

	h.db.Where("id in (?)", scheduleIds).Find(&schedules)

	h.db.Model(&card).Association("Schedules").Replace(schedules)

	h.r.JSON(rw, http.StatusOK, &card)
}

func (h *DBHandler) cardDeleteHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)

	card := Card{}

	h.db.Delete(&card, id)

	h.r.JSON(rw, http.StatusOK, &card)
}
