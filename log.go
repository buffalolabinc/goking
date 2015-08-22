package main

import (
	_ "fmt"
	"net/http"
	"time"
	"database/sql"
)

type Log struct {
	ID int64 `json:"id"`
	Code string `json:"code"`
	Card Card	`json:"card"`
	CardID sql.NullInt64 `json:"-"`
	Schedule Schedule `json:"schedule"`	
	ScheduleID sql.NullInt64 `json:"-"`
	ValidPin bool `json:"valid_pin"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *DBHandler) logsIndexHandler(rw http.ResponseWriter, req *http.Request) {
	page := getPage(req) - 1
	perPage := getPerPage(req)
	offset := perPage * page

	var logs []Log

	h.db.Order("created_at desc").Limit(perPage).Offset(offset).Preload("Card").Preload("Schedule").Find(&logs)

	if logs == nil {
		h.r.JSON(rw, http.StatusOK, make([]int64, 0))
	} else {
		var count int
		h.db.Table("logs").Count(&count)
		vals := make([]interface{}, len(logs))
		for i, v := range logs {
			vals[i] = v
		}

		resp := getResponse(vals, count)

		h.r.JSON(rw, http.StatusOK, resp)
	}
}
