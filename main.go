package main

import (
	"log"
)

func main() {
	c, err := ConfigurationLoad("config.json")
	if err != nil {
		log.Panic(err)
	}

	a := new(App)
	err = a.Initialize(c)
	if err != nil {
		log.Panic(err)
	}

	a.Run(c.GetAddr())
}
