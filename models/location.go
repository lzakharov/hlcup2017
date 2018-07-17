package models

import (
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
func GetLocationAverageMark(id string, params map[string][]string) (LocationAvgMark, error) {
	const age = "date_part('year', age(to_timestamp(users.birth_date)))"
	params["id"] = []string{id}

	conditions := map[string]string{
		"id":       `locations.id=:id`,
		"fromDate": "visits.visited_at>:fromDate",
		"toDate":   "visits.visited_at<:toDate",
		"fromAge":  age + ">:fromAge",
		"toAge":    age + "<:toAge",
		"gender":   "users.gender=:gender"}

	where := []string{}
	args := map[string]interface{}{}

	for param, condition := range conditions {
		if _, ok := params[param]; ok {
			where = append(where, condition)
			args[param] = params[param][0]
		}
	}

	query := `SELECT COALESCE("round"("avg"(visits.mark), 2), 0) as "avg" 
	FROM locations 
	JOIN visits ON visits.location = locations.id 
	JOIN users ON users.id = visits."user" 
	WHERE ` + strings.Join(where, " AND ")

	average := LocationAvgMark{}
	nstmt, err := DB.PrepareNamed(query)
	if err != nil {
		return average, err
	}
	err = nstmt.Get(&average, args)
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
func PopulateLocations(locations Locations) error {
	for _, location := range locations.Rows {
		if err := InsertLocation(location); err != nil {
			return err
		}
	}
	return nil
}
