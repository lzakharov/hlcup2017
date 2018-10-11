package main

// Place contains short description about user-visited location.
type Place struct {
	Mark      uint8  `json:"mark" db:"mark"`
	VisitedAt int32  `json:"visited_at" db:"visited_at"`
	Place     string `json:"place" db:"place"`
}

// Places contains slice of places.
type Places struct {
	Rows []*Place `json:"visits"`
}

// PlaceFilter contains places filtering parameters from requests.
type PlaceFilter struct {
	FromDate *int32  `schema:"fromDate"`
	ToDate   *int32  `schema:"toDate"`
	Country  *string `schema:"country"`
	Distance *uint32 `schema:"distance"`
}
