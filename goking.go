package main

import (
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/asaskevich/govalidator"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mholt/binding"
	_ "github.com/thoas/stats"
	"github.com/unrolled/render"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

type DBHandler struct {
	db *gorm.DB
	r  *render.Render
}

type Model interface {
	GetName() string
}

type AppConfig struct {
	AssetPath string   `json:"asset_path" valid:"required"`
	DbName    string   `json:"db_name" valid:"required"`
	Debug     bool     `json:"debug" valid:"required"`
	DbConfig  []string `json:"db_config" valid:"required"`
	Port      string   `json:"port" valid:"required,numeric"`
}

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

const (
	defaultPerPage = 30
)

func main() {

	// Load configuration
	config := new(AppConfig)
	configPath := GetArgs()
	LoadConfig(configPath, config)
	fmt.Println(skittles.Green("Success:") + " configuration loaded.")

	if config.Debug {
		fmt.Println(skittles.BlinkRed("\n\nWARNING PENDING TABLE TRUNCATE:") + " Do you wish to continue? (y/n)")

		var cont string
		n, err := fmt.Scanf("%s", &cont)
		CheckErr(err)

		if n != 1 || cont == "n" {
			os.Exit(1)
		}
	}

	// setup db
	db, err := gorm.Open("sqlite3", "./"+config.DbName+"?"+strings.Join(config.DbConfig, "?"))
	CheckErr(err)

	db.LogMode(config.Debug)
	defer db.Close()

	models := [...]Model{
		&Schedule{}, &Card{},
	}

	for _, m := range models {
		if config.Debug {
			db.DropTable(m)
		}
		db.CreateTable(m)

		var modelName string = strings.ToLower(strings.Replace(reflect.TypeOf(m).String(), "*main.", "", 1))
		db.Model(m).AddIndex("idx_"+modelName+"_index", "id")
	}

	r := render.New(render.Options{})
	h := DBHandler{db: &db, r: r}

	router := mux.NewRouter()

	router.HandleFunc("/api/schedules", h.schedulesIndexHandler).Methods("GET")

	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("public/web")))

	n.UseHandler(router)
	n.Run(":" + config.Port)

}

func (h *DBHandler) schedulesIndexHandler(rw http.ResponseWriter, req *http.Request) {

}
