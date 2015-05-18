package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"io/ioutil"
	"os"
)

func GetArgs() string {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Printf("Must specify asset location\n\nUsage: %s [asset_path]\n", os.Args[0])
		os.Exit(1)
	}

	return args[0]
}

func LoadConfig(path string, config *AppConfig) {
	dat, err := ioutil.ReadFile(path)
	CheckErr(err)

	jsonErr := json.Unmarshal(dat, config)
	CheckErr(jsonErr)

	_, validErr := govalidator.ValidateStruct(*config)
	CheckErr(validErr)

}

func GetAssetPath(config *AppConfig, path string) string {

	var buffer bytes.Buffer

	buffer.WriteString(config.AssetPath)
	buffer.WriteString("/")
	buffer.WriteString(path)

	return buffer.String()
}

func CheckErr(e error) bool {
	if e != nil {
		panic(e)
	}

	return true
}
