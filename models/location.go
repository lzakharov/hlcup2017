package models

import (
	"fmt"
	"strings"
)

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

// LocationAvgMark contains location average mark.
type LocationAvgMark struct {
	Avg float64 `json:"avg" db:"avg"`
}

// GetLocation returns location from database specified by id.
func GetLocation(id string) (Location, error) {
	location := Location{}
	err := GetByID(locationsTableName, id, &location)
	return location, err
}

// GetLocationAverageMark returns average mark for specified location.
func GetLocationAverageMark(id string, predicates map[string][]string) (LocationAvgMark, error) {
	const age = "date_part('year', age(to_timestamp(users.birth_date)))"
	var (
		names      = []string{"fromDate", "toDate", "fromAge", "toAge", "gender"}
		statements = []string{
			"visits.visited_at > ",
			"visits.visited_at < ",
			age + " > ",
			age + " < ",
			"users.gender = "}
		where   = []string{"locations.id = $1"}
		values  = []interface{}{id}
		average = LocationAvgMark{}
	)

	for i, name := range names {
		if value, ok := predicates[name]; ok {
			values = append(values, value[0])
			where = append(where, fmt.Sprintf("%s$%d", statements[i], len(values)))
		}
	}

	q := `SELECT COALESCE("round"("avg"(visits.mark), 2), 0) as "avg"
		  FROM locations 
		  JOIN visits ON visits.location = locations.id 
		  JOIN users ON users.id = visits."user" 
		  WHERE ` + strings.Join(where, " AND ")

	err := DB.Get(&average, q, values...)
	return average, err
}

// InsertLocation inserts specified location into database.
func InsertLocation(location *Location) error {
	_, err := DB.NamedExec(
		`INSERT INTO locations (id, place, country, city, distance) 
		 VALUES (:id, :place, :country, :city, :distance)`, location)
	return err
}

// PopulateLocations inserts specified list of locations into database.
func PopulateLocations(locations Locations) {
	for _, location := range locations.Rows {
		InsertLocation(location)
	}
}
