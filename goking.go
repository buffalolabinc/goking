package main

import (
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/thoas/stats"
	"github.com/unrolled/render"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	defaultPerPage = 30
)

func main() {

	// Load configuration
	config := new(AppConfig)
	configPath := getArgs()
	loadConfig(configPath, config)
	fmt.Println(skittles.Green("Success:") + " configuration loaded.")

	if config.Debug {
		fmt.Println(skittles.BlinkRed("\n\nWARNING PENDING TABLE TRUNCATE:") + " Do you wish to continue? (y/n)")

		var cont string
		n, err := fmt.Scanf("%s", &cont)
		checkErr(err)

		if n != 1 || cont == "n" {
			os.Exit(1)
		}
	}

	// setup db
	db, err := gorm.Open("sqlite3", "./"+config.DbName+"?"+strings.Join(config.DbConfig, "?"))
	checkErr(err)

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
	router.HandleFunc("/api/logs", h.logsIndexHandler).Methods("GET")
	router.HandleFunc("/api/cards", h.cardsIndexHandler).Methods("GET")

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewStatic(http.Dir("public/web")),
	)

	n.UseHandler(router)
	n.Run(":" + config.Port)
}

func (h *DBHandler) cardsIndexHandler(rw http.ResponseWriter, req *http.Request) {

}

func (h *DBHandler) logsIndexHandler(rw http.ResponseWriter, req *http.Request) {

}
