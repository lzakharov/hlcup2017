package main

// Location contains location database record.
type Location struct {
	ID       *uint32 `json:"id" db:"id"`
	Place    *string `json:"place" db:"place"`
	Country  *string `json:"country" db:"country"`
	City     *string `json:"city" db:"city"`
	Distance *uint32 `json:"distance" db:"distance"`
}

// Locations contains slice of locations.
type Locations struct {
	Rows []*Location `json:"locations"`
}

// LocationFilter contains locations filtering parameters from requests.
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
