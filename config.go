package main

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration represents project configuration data.
type Configuration struct {
	DB struct {
		Host     string `json:"host"`
		Database string `json:"database"`
	} `json:"db"`
	Data string `json:"data"`
	Host string `json:"host"`
	Port string `json:"port"`
}

// LoadConfiguration returns project configuration loaded from the JSON file.
func LoadConfiguration(file string) Configuration {
	var config Configuration

	r, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	dec := json.NewDecoder(r)
	dec.Decode(&config)

	return config
}
