package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func LoadData(archive string, d *Database) error {
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
				return err
			}

			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			options := strings.Split(string(bytes), "\n")
			tm, test := options[0], options[1]
			log.Println("Data generation timestamp", tm, "with type", test)

			f.Close()
			continue
		}

		entity := name[:strings.LastIndex(name, "_")]

		reader, err := file.Open()
		if err != nil {
			return err
		}

		switch entity {
		case "users":
			users := new(Users)
			if err := parse(reader, &users); err != nil {
				return err
			}
			if err := d.PopulateUsers(users); err != nil {
				return err
			}
		case "locations":
			locations := new(Locations)
			if err := parse(reader, &locations); err != nil {
				return err
			}
			if err := d.PopulateLocations(locations); err != nil {
				return err
			}
		case "visits":
			visits := new(Visits)
			if err := parse(reader, &visits); err != nil {
				return err
			}
			if err := d.PopulateVisits(visits); err != nil {
				return err
			}
		}

		reader.Close()
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
