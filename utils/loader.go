package utils

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/lzakharov/hlcup2017/models"
)

// LoadData loads data from archive.
func LoadData(archive string) error {
	log.Println("Loading data from", archive)
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		name := file.Name[:strings.LastIndex(file.Name, ".")]

		if name == "options" {
			f, err := file.Open()
			if err != nil {
				log.Panic(err)
			}
			defer f.Close()

			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				log.Panic(err)
			}

			options := strings.Split(string(bytes), "\n")
			tm, test := options[0], options[1]
			log.Println("Data generation timestamp", tm, "with type", test)

			continue
		}

		entity := name[:strings.LastIndex(name, "_")]

		reader, err := file.Open()
		if err != nil {
			log.Panic(err)
		}
		defer reader.Close()

		switch entity {
		case "users":
			var users models.Users
			if err := parse(reader, &users); err != nil {
				log.Panic(err)
			}
			models.PopulateUsers(users)
		case "locations":
			var locations models.Locations
			if err := parse(reader, &locations); err != nil {
				log.Panic(err)
			}
			models.PopulateLocations(locations)
		case "visits":
			var visits models.Visits
			if err := parse(reader, &visits); err != nil {
				log.Panic(err)
			}
			models.PopulateVisits(visits)
		}
	}

	log.Println("Loaded data from", archive)
	return nil
}

func parse(reader io.ReadCloser, v interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}
