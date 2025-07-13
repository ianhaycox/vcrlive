package model

type LivePositions struct {
	Weekend Weekend  `json:"weekend,omitempty"`
	Session Session  `json:"session,omitempty"`
	Drivers []Driver `json:"drivers,omitempty"`
}
