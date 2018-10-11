package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config contains server properties.
type Config struct {
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	DBConfig *DBConfig `json:"db"`
	Data     string    `json:"data"`
}

// DBConfig contains server database properties.
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

// ConfigurationLoad returns server configurations parsed from the specified file.
func ConfigurationLoad(file string) *Config {
	config := new(Config)

	r, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer r.Close()

	dec := json.NewDecoder(r)
	if err := dec.Decode(config); err != nil {
		log.Panic(err)
	}

	return config
}

// GetDSN returns database source name.
func (d *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

// GetAddr returns sever address.
func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
