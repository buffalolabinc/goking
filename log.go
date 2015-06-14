package main

import (
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type Log struct {
	ID        int64     `json:"Id"`
	Code      string    `json:"Code"`
	ValidPin  bool      `json:"ValidPin"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func (l *Log) GetName() string {
	return "log"
}

type LogForm struct {
	Code     string
	ValidPin bool
}

func (lf *LogForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&lf.Code: binding.Field{
			Form:     "code",
			Required: true,
		},
		&lf.Code: binding.Field{
			Form:     "validpin",
			Required: true,
		},
	}
}

func (h *DBHandler) logsIndexHandler(rw http.ResponseWriter, req *http.Request) {
	page := getPage(req) - 1
	perPage := getPerPage(req)
	offset := perPage * page

	var logs []Log

	h.db.Limit(perPage).Offset(offset).Find(&logs)

	if logs == nil {
		h.r.JSON(rw, http.StatusOK, make([]int64, 0))
	} else {
		h.r.JSON(rw, http.StatusOK, &logs)
	}
}
