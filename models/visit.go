package models

import "log"

// Visit contains information about visit.
type Visit struct {
	ID        uint32 `json:"id" db:"id"`
	Location  uint32 `json:"location" db:"location"`
	User      uint32 `json:"user" db:"user"`
	VisitedAt int32  `json:"visited_at" db:"visited_at"`
	Mark      uint8  `json:"mark" db:"mark"`
}

// Visits contains list of visits.
type Visits struct {
	Rows []*Visit `json:"visits"`
}

// GetVisit returns visit from database specified by id.
func GetVisit(id uint32) (Visit, error) {
	visit := Visit{}
	if err := DB.Get(&visit, "SELECT * FROM visits WHERE id=$1", id); err != nil {
		return visit, err
	}
	return visit, nil
}

// InsertVisit inserts specified visit into database.
func InsertVisit(visit *Visit) {
	_, err := DB.NamedExec(
		`INSERT INTO visits (id, location, "user", visited_at, mark) 
		 VALUES (:id, :location, :user, :visited_at, :mark)`, visit)
	if err != nil {
		log.Fatal(err)
	}
}

// PopulateVisits inserts specified list of Visits into database.
func PopulateVisits(visits Visits) {
	for _, visit := range visits.Rows {
		InsertVisit(visit)
	}
}
