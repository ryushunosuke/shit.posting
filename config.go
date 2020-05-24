package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var config Config

// Config is the struct used to reading from the config.json
type Config struct {
	DB              string   `json:"DB"`
	Folders         []string `json:"Folders"`
	Port            string   `json:"Port"`
	Types           []string `json:"FileTypes"`
	TypeMap         map[string]bool
	ThumbnailFolder string `json:"Thumbnail"`
	StrFilesize     string `json:"Filesize"`
	Filesize        int64
}

// LoadConfig loads ./config.json into the Config struct
func LoadConfig() (Config, error) {
	// var config Config
	// Default config settings
	config.TypeMap = map[string]bool{}
	config.ThumbnailFolder = "./thumbnail/"
	config.Port = "8000"
	config.Filesize = 1024 * 1024 * 1024

	file, err := os.Open("./config.json")
	if err != nil {
		fmt.Print(err.Error())
		return config, err
	}
	defer file.Close()
	JSONParser := json.NewDecoder(file)
	err = JSONParser.Decode(&config)
	if err != nil {
		fmt.Print(err.Error())
		return config, err
	}
	if config.StrFilesize != "" {
		config.Filesize = StringToInt(config.StrFilesize)
	}
	for _, val := range config.Types {
		config.TypeMap[val] = true
	}
	return config, err
}

// StringToInt handles strings such as "100GB" from the config file and converts them to bytes.
func StringToInt(size string) (Size int64) {
	suffix := ""
	if len(size) > 2 {
		suffix = size[len(size)-2:]
		size = size[:len(size)-2]
	}
	Size, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		fmt.Println("Error parsing Filesize config.")
		return 0
	}
	switch suffix {
	case "GB":
		Size *= 1024 * 1024 * 1024
	case "MB":
		Size *= 1024 * 1024
	case "KB":
		Size *= 1024

	}
	return

}
