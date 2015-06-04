package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type Schedule struct {
	ID        int64     `json:"Id"`
	Name      string    `json:"Name"`
	Mon       bool      `json:"Mon"`
	Tue       bool      `json:"Tue"`
	Wed       bool      `json:"Wed`
	Thu       bool      `json:"Thu"`
	Fri       bool      `json:"Fri"`
	Sat       bool      `json:"Sat"`
	Sun       bool      `json:"Sun"`
	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	DeletedAt time.Time `json:"DeletedAt"`
}

func (s *Schedule) GetName() string {
	return "schedule"
}

type ScheduleForm struct {
	Name      string    `valid:"alpha,required"`
	Mon       bool      `valid:"required"`
	Tue       bool      `valid:"required"`
	Wed       bool      `valid:"required"`
	Thu       bool      `valid:"required"`
	Fri       bool      `valid:"required"`
	Sat       bool      `valid:"required"`
	Sun       bool      `valid:"required"`
	StartTime time.Time `valid:"required"`
	EndTime   time.Time `valid:"required"`
}

func (sf *ScheduleForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&sf.Name: binding.Field{
			Form:     "Name",
			Required: true,
		},
		&sf.Mon: binding.Field{
			Form:     "Mon",
			Required: false,
		},
		&sf.Tue: binding.Field{
			Form:     "Tue",
			Required: false,
		},
		&sf.Wed: binding.Field{
			Form:     "Wed",
			Required: false,
		},
		&sf.Thu: binding.Field{
			Form:     "Thu",
			Required: false,
		},
		&sf.Fri: binding.Field{
			Form:     "Fri",
			Required: false,
		},
		&sf.Sat: binding.Field{
			Form:     "Sat",
			Required: false,
		},
		&sf.Sun: binding.Field{
			Form:     "Sun",
			Required: false,
		},
		&sf.StartTime: binding.Field{
			Form:     "StartTime",
			Required: true,
		},
		&sf.EndTime: binding.Field{
			Form:     "EndTime",
			Required: true,
		},
	}
}

func (sf *ScheduleForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	_, err := govalidator.ValidateStruct(sf)
	if err != nil {
		// validate date start and end / valid times etc
	}
	return errs
}

func (h *DBHandler) schedulesIndexHandler(rw http.ResponseWriter, req *http.Request) {
	page := getPage(req) - 1
	perPage := getPerPage(req)
	offset := perPage * page

	var schedules []Schedule

	h.db.Limit(perPage).Offset(offset).Find(&schedules)

	if schedules == nil {
		h.r.JSON(rw, http.StatusOK, make([]int64, 0))
	} else {
		h.r.JSON(rw, http.StatusOK, &schedules)
	}
}

func (h *DBHandler) scheduleShowHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	schedule := Schedule{}
	h.db.First(&schedule, id)
	h.r.JSON(rw, http.StatusOK, &schedule)
}

func (h *DBHandler) scheduleCreateHandler(rw http.ResponseWriter, req *http.Request) {
	h.schedulesEdit(rw, req, 0)
}

func (h *DBHandler) scheduleUpdateHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	h.schedulesEdit(rw, req, id)
}

func (h *DBHandler) scheduleDeleteHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	schedule := Schedule{}
	h.db.Delete(&schedule, id)
	h.r.JSON(rw, http.StatusOK, &schedule)
}

func (h *DBHandler) schedulesEdit(rw http.ResponseWriter, req *http.Request, id int64) {
	scheduleForm := ScheduleForm{}

	if err := binding.Bind(req, &scheduleForm); err.Handle(rw) {
		return
	}

	fmt.Println("%v+", scheduleForm)

	schedule := Schedule{
		ID:        id,
		Name:      scheduleForm.Name,
		Mon:       scheduleForm.Mon,
		Tue:       scheduleForm.Tue,
		Wed:       scheduleForm.Wed,
		Thu:       scheduleForm.Thu,
		Fri:       scheduleForm.Fri,
		Sat:       scheduleForm.Sat,
		Sun:       scheduleForm.Sun,
		StartTime: scheduleForm.StartTime,
		EndTime:   scheduleForm.EndTime,
	}

	h.db.Save(&schedule)
	h.r.JSON(rw, http.StatusOK, &schedule)
}
