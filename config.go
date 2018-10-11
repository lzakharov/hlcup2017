package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	DBConfig *DBConfig `json:"db"`
	Data     string    `json:"data"`
}

type DBConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"sslmode"`
	Schema   string `json:"schema"`
}

func ConfigurationLoad(file string) *Config {
	config := new(Config)

	r, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	dec := json.NewDecoder(r)
	dec.Decode(config)

	return config
}

func (d *DBConfig) GetDataSourceName() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
