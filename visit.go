package main

type Visit struct {
	ID        *uint32 `json:"id" db:"id"`
	Location  *uint32 `json:"location" db:"location"`
	User      *uint32 `json:"user" db:"user"`
	VisitedAt *int32  `json:"visited_at" db:"visited_at"`
	Mark      *uint8  `json:"mark" db:"mark"`
}

type Visits struct {
	Rows []*Visit `json:"visits"`
}
