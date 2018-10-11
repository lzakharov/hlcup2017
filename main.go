package main

import (
	"log"
)

func main() {
	c := ConfigurationLoad("config.json")

	a := new(App)
	err := a.Initialize(c)
	if err != nil {
		log.Panic(err)
	}

	a.Run(c.GetAddr())
}
