package model

import "encoding/json"

type LivePositions struct {
	Weekend Weekend  `json:"weekend"`
	Session Session  `json:"session"`
	Drivers []Driver `json:"drivers,omitempty"`
}

func (l *LivePositions) String() string {
	b, _ := json.MarshalIndent(l, "", "  ")

	return string(b)
}
