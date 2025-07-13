package model

type LivePositions struct {
	Weekend Weekend  `json:"weekend"`
	Session Session  `json:"session"`
	Drivers []Driver `json:"drivers,omitempty"`
}
