package models

import "log"

// Location contains information about location.
type Location struct {
	ID       uint32 `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance uint32 `json:"distance"`
}

// Locations contains list of locations.
type Locations struct {
	Rows []*Location `json:"locations"`
}

// InsertLocation inserts specified location into database.
func InsertLocation(location *Location) {
	_, err := DB.NamedExec(
		`INSERT INTO locations (id, place, country, city, distance) 
		 VALUES (:id, :place, :country, :city, :distance)`, location)
	if err != nil {
		log.Fatal(err)
	}
}

// PopulateLocations inserts specified list of locations into database.
func PopulateLocations(locations Locations) {
	for _, location := range locations.Rows {
		InsertLocation(location)
	}
}
