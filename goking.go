package main

import (
	_ "bytes"
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

	models := [...]interface{}{
		&Schedule{}, &Card{}, &Log{},
	}

	if config.Truncate {
		db.Exec("DROP TABLE IF EXISTS card_schedule")
	}

	for _, m := range models {
		if config.Truncate {
			db.DropTable(m)
		}
		db.AutoMigrate(m)
	}

	go func(db gorm.DB) {
		c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600, ReadTimeout: time.Second * 5}

		s, err := serial.OpenPort(c)
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1024)
		for true {
			n, _ := s.Read(buf)

			s := string(buf[:n])
			if strings.Contains(s, "\r\n") {
				segments := strings.Split(s, "\r\n")

				_ = "breakpoint"
				// process each segment
				//@TODO This was being strang
				for _, e := range segments {
					subSegs := strings.Split(e, ":")
					fmt.Print(len(subSegs))
					if subSegs[0] != "A" || len(subSegs) != 3 {
						continue
					}

					l := len(subSegs)
					_ = "breakpoint"
					cmd, door_card, pin := subSegs[0], subSegs[1], subSegs[2]
					_ = "breakpoint"
					fmt.Print(cmd, door_card, pin, l)
				}
			}
		}
	}(db)

	configureHttpAndListen(config, db)

}

func configureHttpAndListen(config *AppConfig, db gorm.DB) {
	// register routes
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
