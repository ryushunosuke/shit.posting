package main

import (
	"encoding/json"
	"fmt"
	"os"
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
}

// LoadConfig loads ./config.json into the Config struct
func LoadConfig() (Config, error) {
	// var config Config
	// Default config settings
	config.TypeMap = map[string]bool{}
	config.ThumbnailFolder = "./thumbnail/"
	config.Port = "8000"

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
	for _, val := range config.Types {
		config.TypeMap[val] = true
	}
	return config, err
}
