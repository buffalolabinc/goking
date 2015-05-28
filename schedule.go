package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type Schedule struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Mon       bool      `json:"mon"`
	Tue       bool      `json:"tue"`
	Wed       bool      `json:"wed`
	Thu       bool      `json:"thu"`
	Fri       bool      `json:"fri"`
	Sat       bool      `json:"sat"`
	Sun       bool      `json:"sun"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
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

func (sf *ScheduleForm) FeildMap() binding.FieldMap {
	return binding.FieldMap{
		&sf.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
		&sf.Mon: binding.Field{
			Form:     "mon",
			Required: true,
		},
		&sf.Tue: binding.Field{
			Form:     "tue",
			Required: true,
		},
		&sf.Wed: binding.Field{
			Form:     "wed",
			Required: true,
		},
		&sf.Thu: binding.Field{
			Form:     "thu",
			Required: true,
		},
		&sf.Fri: binding.Field{
			Form:     "fri",
			Required: true,
		},
		&sf.Sat: binding.Field{
			Form:     "sat",
			Required: true,
		},
		&sf.Sun: binding.Field{
			Form:     "sun",
			Required: true,
		},
		&sf.StartTime: binding.Field{
			Form:     "start_time",
			Required: true,
		},
		&sf.EndTime: binding.Field{
			Form:     "end_time",
			Required: true,
		},
	}
}

func (sf ScheduleForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	_, err := govalidator.ValidateStruct(sf)
	if err != nil {
		// do something
	}
	return errs
}
