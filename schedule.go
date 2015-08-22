package main

import (
	"github.com/asaskevich/govalidator"
	"net/http"
	"time"
	"encoding/json"
)

type ScheduleIdentity struct {
	ID	int64	`json:"id" valid:"required"`
}

type Schedule struct {
	ScheduleIdentity
	Name      string    `json:"name"`
	Mon       bool      `json:"mon"`
	Tue       bool      `json:"tue"`
	Wed       bool      `json:"wed"`
	Thu       bool      `json:"thu"`
	Fri       bool      `json:"fri"`
	Sat       bool      `json:"sat"`
	Sun       bool      `json:"sun"`
	StartTime time.Time `json:"start_time" valid:"required"`
	EndTime   time.Time `json:"end_time" valid:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (schedule *Schedule) MatchesNow() bool {
	nowTime := time.Now()
	normalizedTime := time.Date(1970, 1, 1, nowTime.Hour(), nowTime.Minute(), nowTime.Second(), nowTime.Nanosecond(), nowTime.Location())

	if schedule.StartTime.Before(normalizedTime) && schedule.EndTime.After(normalizedTime) {
		switch nowTime.Weekday() {
			case time.Monday:
				return schedule.Mon == true
			case time.Tuesday:
				return schedule.Tue == true
			case time.Wednesday:
				return schedule.Wed == true
			case time.Thursday:
				return schedule.Thu == true
			case time.Friday:
				return schedule.Fri == true
			case time.Saturday:
				return schedule.Sat == true
			case time.Sunday:
				return schedule.Sun == true
		}
	}

	return false
}

type ScheduleRequest struct {
	Name      string    `json:"name"`
	Mon       bool      `json:"mon"`
	Tue       bool      `json:"tue"`
	Wed       bool      `json:"wed"`
	Thu       bool      `json:"thu"`
	Fri       bool      `json:"fri"`
	Sat       bool      `json:"sat"`
	Sun       bool      `json:"sun"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
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
		var count int
		h.db.Table("schedules").Count(&count)

		vals := make([]interface{}, len(schedules))
		for i, v := range schedules {
			vals[i] = v
		}

		resp := getResponse(vals, count)

		h.r.JSON(rw, http.StatusOK, &resp)
	}
}

func (h *DBHandler) scheduleShowHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)

	schedule := Schedule{}

	h.db.First(&schedule, id)

	h.r.JSON(rw, http.StatusOK, &schedule)
}

func (h *DBHandler) scheduleCreateHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	createSchedule := ScheduleRequest{}

	err := decoder.Decode(&createSchedule)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	schedule := Schedule{}

	schedule.Name = createSchedule.Name
	schedule.Mon = createSchedule.Mon
	schedule.Tue = createSchedule.Tue
	schedule.Wed = createSchedule.Wed
	schedule.Thu = createSchedule.Thu
	schedule.Fri = createSchedule.Fri
	schedule.Sat = createSchedule.Sat
	schedule.Sun = createSchedule.Sun
	schedule.StartTime = createSchedule.StartTime
	schedule.EndTime = createSchedule.EndTime

	_, err = govalidator.ValidateStruct(&schedule)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	h.db.Save(&schedule)

	h.r.JSON(rw, http.StatusOK, &schedule)
}

func (h *DBHandler) scheduleUpdateHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)

	decoder := json.NewDecoder(req.Body)

	updateSchedule := ScheduleRequest{}

	err := decoder.Decode(&updateSchedule)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	schedule := Schedule{}

	h.db.First(&schedule, id)

	schedule.Name = updateSchedule.Name
	schedule.Mon = updateSchedule.Mon
	schedule.Tue = updateSchedule.Tue
	schedule.Wed = updateSchedule.Wed
	schedule.Thu = updateSchedule.Thu
	schedule.Fri = updateSchedule.Fri
	schedule.Sat = updateSchedule.Sat
	schedule.Sun = updateSchedule.Sun
	schedule.StartTime = updateSchedule.StartTime
	schedule.EndTime = updateSchedule.EndTime

	_, err = govalidator.ValidateStruct(&schedule)

	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error":err.Error()})
		return
	}

	h.db.Save(&schedule)

	h.r.JSON(rw, http.StatusOK, &schedule)
}

func (h *DBHandler) scheduleDeleteHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)

	schedule := Schedule{}

	h.db.Delete(&schedule, id)

	h.r.JSON(rw, http.StatusOK, &schedule)
}
