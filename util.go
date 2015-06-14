package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type AppConfig struct {
	AssetPath string   `json:"asset_path" valid:"required"`
	DbName    string   `json:"db_name" valid:"required"`
	Debug     bool     `json:"debug" valid:"required"`
	DbConfig  []string `json:"db_config" valid:"required"`
	Port      string   `json:"port" valid:"required,numeric"`
	Truncate  bool     `json:"truncate"`
}

func getArgs() string {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Printf("Must specify asset location\n\nUsage: %s [asset_path]\n", os.Args[0])
		os.Exit(1)
	}

	return args[0]
}

func loadConfig(path string, config *AppConfig) {
	dat, err := ioutil.ReadFile(path)
	checkErr(err)

	jsonErr := json.Unmarshal(dat, config)
	checkErr(jsonErr)

	_, validErr := govalidator.ValidateStruct(*config)
	checkErr(validErr)

}

func getAssetPath(config *AppConfig, path string) string {

	var buffer bytes.Buffer

	buffer.WriteString(config.AssetPath)
	buffer.WriteString("/")
	buffer.WriteString(path)

	return buffer.String()
}

func checkErr(e error) bool {
	if e != nil {
		panic(e)
	}

	return true
}

func getResponse(obj []interface{}) *PaginatedResponse {
	return &PaginatedResponse{Items: obj, TotalItems: len(obj)}
}

func getId(req *http.Request) int64 {
	vars := mux.Vars(req)
	idString := vars["id"]
	id, err := strconv.ParseInt(idString, 10, 0)
	if err != nil {
		log.Println(skittles.BoldRed(err))
	}

	return id
}

func getPage(req *http.Request) int {
	return parseQueryValues(req, "page")
}

func getPerPage(req *http.Request) int {
	perPage := parseQueryValues(req, "per_page")

	if perPage == 0 {
		return defaultPerPage
	}

	return perPage
}

func parseQueryValues(req *http.Request, value string) int {
	vals := req.URL.Query()
	val := vals[value]
	if val != nil {
		v, err := strconv.ParseInt(val[0], 10, 0)

		if err != nil {
			log.Println(skittles.BoldRed(err))
		}

		return int(v)
	}

	return 0
}
