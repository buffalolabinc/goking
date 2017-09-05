package main

import (
	_ "bytes"
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/tarm/serial"
	_ "github.com/thoas/stats"
	"github.com/unrolled/render"
	"os"
	"strings"
	"log"
	"time"
	"bufio"
	"encoding/binary"
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
	fmt.Println("Connecting with %s and DSN %s", config.DbDriver, config.DbDsn)
  db, err := gorm.Open(config.DbDriver, config.DbDsn)
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

	if config.Serial.Enabled {
		fmt.Println("Starting Serial goroutine")

		go func(db *gorm.DB) {
			c := &serial.Config{Name: config.Serial.DevicePath, Baud: config.Serial.BaudRate, ReadTimeout: time.Second * 5}

			s, err := serial.OpenPort(c)
			if err != nil {
				log.Fatal(err)
			}

			reader := bufio.NewReader(s)
			writer := bufio.NewWriter(s)
	SerialLoop:
			for true {
				reply, _, err := reader.ReadLine()

				if err != nil {
					panic(err)
				}

				command := string(reply[:])

				subSegs := strings.Split(command, ":")

				fmt.Println("number of segments ", len(subSegs))

				if subSegs[0] != "A" || len(subSegs) != 3 {
					continue
				}

				l := len(subSegs)

				cmd, door_card, pin := subSegs[0], subSegs[1], subSegs[2]

				log := Log{
					Code: door_card,
				}

				fmt.Println("cmd ", cmd, " door card ", door_card, " pin ", pin, l)

				card := Card{}

				db.Where(&Card{Code: door_card, IsActive: true}).First(&card)

				log.Card = card

				if card.Pin != pin {
					fmt.Println("Pin does not match")

					log.ValidPin = false
					db.Save(&log)

					continue SerialLoop
				}

				log.ValidPin = true

				db.Model(&card).Association("Schedules").Find(&card.Schedules)

				for _, schedule := range card.Schedules {
					if schedule.MatchesNow() {
						fmt.Println("Matching schedule found! ", schedule.Name)
						// send open door

						err := binary.Write(writer, binary.BigEndian, byte(0))

						if err != nil {
							fmt.Println("Failed to send response over Serial ", err)
							continue SerialLoop
						}

						err = binary.Write(writer, binary.BigEndian, uint32(5))

						if err != nil {
							fmt.Println("Failed to send response over Serial ", err)
							continue SerialLoop
						}

						log.Schedule = schedule
						db.Save(&log)

						continue SerialLoop
					}
				}


			}
		}(db)
	}

	configureHttpAndListen(config, db)

}

func configureHttpAndListen(config *AppConfig, db *gorm.DB) {
	// register routes
	r := render.New(render.Options{})
	h := DBHandler{db: db, r: r}

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
	)

	n.UseHandler(router)
	n.Run(":" + config.Port)
}
