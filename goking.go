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
		fmt.Println(skittles.Green("\nDebug mode enabled"))
	}

	if config.Truncate {
		fmt.Println(skittles.BlinkRed("\n\nWARNING PENDING TABLE TRUNCATE:") + " Do you wish to continue? (y/n)")

		var cont string
		n, err := fmt.Scanf("%s", &cont)
		checkErr(err)

		if n != 1 || cont == "n" {
			os.Exit(1)
		}
	}

	// setup db
	fmt.Println("%v+", config.DbConfig)
	db, err := gorm.Open("sqlite3", "./"+config.DbName+"?"+strings.Join(config.DbConfig, "?"))
	checkErr(err)

	db.LogMode(config.Debug)
	defer db.Close()

	models := [...]Model{
		&Schedule{}, &Card{},
	}

	for _, m := range models {
		if config.Truncate {
			db.DropTable(m)
		}
	}

	db.AutoMigrate(&Card{}, &Schedule{})

	r := render.New(render.Options{})
	h := DBHandler{db: &db, r: r}

	router := mux.NewRouter()

	router.HandleFunc("/api/schedules", h.schedulesIndexHandler).Methods("GET")
	router.HandleFunc("/api/schedules", h.scheduleCreateHandler).Methods("POST")
	router.HandleFunc("/api/schedules/{id:[0-9]+}", h.scheduleShowHandler).Methods("GET")
	router.HandleFunc("/api/schedules/{id:[0-9]+}", h.scheduleUpdateHandler).Methods("PUT", "PATCH")
	router.HandleFunc("/api/schedules/{id:[0-9]+}", h.scheduleDeleteHandler).Methods("DELETE")

	router.HandleFunc("/api/cards", h.cardsIndexHandler).Methods("GET")
	router.HandleFunc("/api/cards", h.cardCreateHandler).Methods("POST")
	router.HandleFunc("/api/cards/{id:[0-9]+}", h.cardShowHandler).Methods("GET")
	router.HandleFunc("/api/cards/{id:[0-9]+}", h.cardUpdateHandler).Methods("PUT", "PATCH")
	router.HandleFunc("/api/cards/{id:[0-9]+}", h.cardDeleteHandler).Methods("DELETE")

	router.HandleFunc("/api/logs", h.logsIndexHandler).Methods("GET")

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewStatic(http.Dir("public/web")),
	)

	n.UseHandler(router)
	n.Run(":" + config.Port)
}

func (h *DBHandler) logsIndexHandler(rw http.ResponseWriter, req *http.Request) {

}
