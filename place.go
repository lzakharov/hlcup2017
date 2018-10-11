package main

type Place struct {
	Mark      uint8  `json:"mark" db:"mark"`
	VisitedAt int32  `json:"visited_at" db:"visited_at"`
	Place     string `json:"place" db:"place"`
}

type Places struct {
	Rows []*Place `json:"visits"`
}

type PlaceFilter struct {
	FromDate *int32  `schema:"fromDate"`
	ToDate   *int32  `schema:"toDate"`
	Country  *string `schema:"country"`
	Distance *uint32 `schema:"distance"`
}
