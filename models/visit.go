package models

const visitsTableName = "visits"

// Visit contains full information about visit.
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
func GetVisit(id string) (Visit, error) {
	visit := Visit{}
	err := GetByID(visitsTableName, id, &visit)
	return visit, err
}

// InsertVisit inserts specified visit into database.
func InsertVisit(visit *Visit) error {
	_, err := DB.NamedExec(
		`INSERT INTO visits (id, location, "user", visited_at, mark) 
		 VALUES (:id, :location, :user, :visited_at, :mark)`, visit)
	return err
}

// PopulateVisits inserts specified list of Visits into database.
func PopulateVisits(visits Visits) error {
	for _, visit := range visits.Rows {
		if err := InsertVisit(visit); err != nil {
			return err
		}
	}
	return nil
}
