package models

import (
	"log"

	sq "github.com/Masterminds/squirrel"
)

const visitsTableName = "visits"

var visitsTableColumns = []string{"id", "location", `"user"`, "visited_at", "mark"}

// Visit contains full information about visit.
type Visit struct {
	ID        *uint32 `json:"id" db:"id"`
	Location  *uint32 `json:"location" db:"location"`
	User      *uint32 `json:"user" db:"user"`
	VisitedAt *int32  `json:"visited_at" db:"visited_at"`
	Mark      *uint8  `json:"mark" db:"mark"`
}

// Visits contains list of visits.
type Visits struct {
	Rows []*Visit `json:"visits"`
}

// GetVisit returns visit from database specified by id.
func GetVisit(id string) (*Visit, error) {
	visit := new(Visit)
	err := GetByID(visitsTableName, id, visit)
	return visit, err
}

// InsertVisit inserts specified visit into database.
func InsertVisit(visit *Visit) error {
	sql, args, err := psql.
		Insert(visitsTableName).
		Columns(visitsTableColumns...).
		Values(visit.ID, visit.Location, visit.User, visit.VisitedAt, visit.Mark).
		ToSql()

	_, err = DB.Exec(sql, args...)
	if err != nil {
		log.Println(err)
	}
	return err
}

// PopulateVisits inserts specified list of Visits into database.
func PopulateVisits(visits *Visits) error {
	for _, visit := range visits.Rows {
		if err := InsertVisit(visit); err != nil {
			return err
		}
	}
	return nil
}

// UpdateVisit updates specified visit's row in database.
func UpdateVisit(id string, visit *Visit) error {
	update := psql.Update(visitsTableName)

	if visit.Location != nil {
		update = update.Set("location", visit.Location)
	}
	if visit.User != nil {
		update = update.Set(`"user"`, visit.User)
	}
	if visit.VisitedAt != nil {
		update = update.Set("visited_at", visit.VisitedAt)
	}
	if visit.Mark != nil {
		update = update.Set("mark", visit.Mark)
	}

	update = update.Where(sq.Eq{"id": id})

	sql, args, err := update.ToSql()

	_, err = DB.Exec(sql, args...)
	if err != nil {
		log.Println(err)
	}

	return err
}
