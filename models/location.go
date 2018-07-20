package models

import (
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
)

const locationsTableName = "locations"

var locationsTableColumns = []string{"id", "place", "country", "city", "distance"}

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

// LocationFilter contains filtering parameters for locations.
type LocationFilter struct {
	FromDate *int32  `schema:"fromDate"`
	ToDate   *int32  `schema:"toDate"`
	FromAge  *int32  `schema:"fromAge"`
	ToAge    *int32  `schema:"toAge"`
	Gender   *string `schema:"gender"`
}

// LocationAvgMark contains location average mark.
type LocationAvgMark struct {
	Avg float64 `json:"avg" db:"avg"`
}

// GetLocation returns location from database specified by id.
func GetLocation(id string) (*Location, error) {
	location := new(Location)
	err := GetByID(locationsTableName, id, location)
	return location, err
}

// GetLocationAverageMark returns average mark for specified location.
func GetLocationAverageMark(id string, filter *LocationFilter) (*LocationAvgMark, error) {
	locations := psql.
		Select(`COALESCE("round"("avg"(visits.mark), 2), 0) AS "avg"`).
		From(locationsTableName).
		Join(fmt.Sprintf("%s ON %s.id = %s.location", visitsTableName, locationsTableName, visitsTableName)).
		Join(fmt.Sprintf(`%s ON %s."user" = %s.id`, usersTableName, visitsTableName, usersTableName)).
		Where(sq.Eq{locationsTableName + ".id": id})

	if filter.FromDate != nil {
		locations = locations.Where(sq.Gt{"visits.visited_at": *filter.FromDate})
	}
	if filter.ToDate != nil {
		locations = locations.Where(sq.Lt{"visits.visited_at": *filter.ToDate})
	}

	const age = "date_part('year', age(to_timestamp(users.birth_date)))"
	if filter.FromAge != nil {
		locations = locations.Where(sq.Gt{age: *filter.FromAge})
	}
	if filter.ToAge != nil {
		locations = locations.Where(sq.Lt{age: *filter.ToAge})
	}

	if filter.Gender != nil {
		locations = locations.Where(sq.Eq{"users.gender": filter.Gender})
	}

	sql, args, err := locations.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	average := &LocationAvgMark{}
	if err := DB.Get(average, sql, args...); err != nil {
		log.Println(err)
		return nil, err
	}

	return average, err
}

// InsertLocation inserts specified location into database.
func InsertLocation(location *Location) error {
	sql, args, err := psql.
		Insert(locationsTableName).
		Columns(locationsTableColumns...).
		Values(location.ID, location.Place, location.Country, location.City, location.Distance).
		ToSql()

	_, err = DB.Exec(sql, args...)
	if err != nil {
		log.Println(err)
	}
	return err
}

// PopulateLocations inserts specified list of locations into database.
func PopulateLocations(locations *Locations) error {
	for _, location := range locations.Rows {
		if err := InsertLocation(location); err != nil {
			return err
		}
	}
	return nil
}

// UpdateLocation updates specified location's row in database.
func UpdateLocation(params map[string]interface{}) error {
	query := prepareUpdate(usersTableName, []string{"place", "country", "city", "distance"}, params)
	_, err := DB.NamedExec(query, params)
	return err
}
