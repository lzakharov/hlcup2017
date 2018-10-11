package main

type Location struct {
	ID       *uint32 `json:"id" db:"id"`
	Place    *string `json:"place" db:"place"`
	Country  *string `json:"country" db:"country"`
	City     *string `json:"city" db:"city"`
	Distance *uint32 `json:"distance" db:"distance"`
}

type Locations struct {
	Rows []*Location `json:"locations"`
}

type LocationFilter struct {
	FromDate *int32  `schema:"fromDate"`
	ToDate   *int32  `schema:"toDate"`
	FromAge  *int32  `schema:"fromAge"`
	ToAge    *int32  `schema:"toAge"`
	Gender   *string `schema:"gender"`
}

type LocationAvgMark struct {
	Avg float64 `json:"avg" db:"avg"`
}
