package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// LoadData loads users, locations and visits from the specified archive to the database.
func LoadData(archive string, d *Database) error {
	log.Println("Loading data from", archive)
	zipReader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
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

			if err = f.Close(); err != nil {
				return err
			}
			continue
		}

		entity := name[:strings.LastIndex(name, "_")]
		reader, err := file.Open()
		if err != nil {
			return err
		}
		if err = loadEntity(entity, &reader, d); err != nil {
			return err
		}

		if err = reader.Close(); err != nil {
			return err
		}
	}

	log.Println("Loaded data from", archive)
	return nil
}

func parse(reader *io.ReadCloser, v interface{}) error {
	data, err := ioutil.ReadAll(*reader)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}

func loadEntity(entity string, r *io.ReadCloser, d *Database) error {
	switch entity {
	case "users":
		if err := loadUsers(r, d); err != nil {
			return err
		}
	case "locations":
		if err := loadLocations(r, d); err != nil {
			return err
		}
	case "visits":
		if err := loadVisits(r, d); err != nil {
			return err
		}
	}
	return nil
}

func loadUsers(r *io.ReadCloser, d *Database) error {
	users := new(Users)
	if err := parse(r, &users); err != nil {
		return err
	}
	return d.PopulateUsers(users)
}

func loadLocations(r *io.ReadCloser, d *Database) error {
	locations := new(Locations)
	if err := parse(r, &locations); err != nil {
		return err
	}
	return d.PopulateLocations(locations)
}

func loadVisits(r *io.ReadCloser, d *Database) error {
	visits := new(Visits)
	if err := parse(r, &visits); err != nil {
		return err
	}
	return d.PopulateVisits(visits)
}
