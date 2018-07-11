package models

// Visit contains information about visit.
type Visit struct {
	ID        uint32 `json:"id"`
	Location  uint32 `json:"location"`
	User      uint32 `json:"user"`
	VisitedAt int32  `json:"visited_at"`
	Mark      uint8  `json:"mark"`
}
