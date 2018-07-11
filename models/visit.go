package models

import "log"

// Visit contains information about visit.
type Visit struct {
	ID        uint32 `json:"id"`
	Location  uint32 `json:"location"`
	User      uint32 `json:"user"`
	VisitedAt int32  `json:"visited_at"`
	Mark      uint8  `json:"mark"`
}

// Visits contains list of visits.
type Visits struct {
	Rows []*Visit `json:"visits"`
}

// InsertVisit inserts specified visit into database.
func InsertVisit(visit *Visit) {
	_, err := DB.NamedExec(
		`INSERT INTO visits (id, location, "user", visited_at, mark) 
		 VALUES (:id, :location, :user, :visitedat, :mark)`, visit)
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
