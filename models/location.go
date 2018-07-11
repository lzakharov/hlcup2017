package models

import "log"

const locationsTableName = "locations"

// Location contains information about location.
type Location struct {
	ID       uint32 `json:"id" db:"id"`
	Place    string `json:"place" db:"place"`
	Country  string `json:"country" db:"country"`
	City     string `json:"city" db:"city"`
	Distance uint32 `json:"distance" db:"distance"`
}

// Locations contains list of locations.
type Locations struct {
	Rows []*Location `json:"locations"`
}

// GetLocation returns location from database specified by id.
func GetLocation(id string) (Location, error) {
	location := Location{}
	if err := GetByID(locationsTableName, id, &location); err != nil {
		return location, err
	}
	return location, nil
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
