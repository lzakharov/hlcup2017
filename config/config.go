package config

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration represents project configuration data.
type Configuration struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	DB   struct {
		Driver   string `json:"driver"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		SSLMode  string `json:"sslmode"`
		Schema   string `json:"schema"`
	} `json:"db"`
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
