package model

import "github.com/ianhaycox/vcrlive/irsdk/iryaml"

const (
	invalid = iota
	getInCar
	warmup
	paradeLaps
	racing
	checkered
	coolDown
)

type Session struct {
	SessionNum   int    `json:"session_num,omitempty"`
	SessionLaps  string `json:"session_laps,omitempty"`
	SessionType  string `json:"session_type,omitempty"`
	SessionName  string `json:"session_name,omitempty"`
	SessionState string `json:"session_state,omitempty"`
}

func NewSession(sessionNum int, sessions []iryaml.Session) Session {
	session := Session{}

	for i := range sessions {
		if sessionNum == sessions[i].SessionNum {
			return Session{
				SessionNum:  sessions[i].SessionNum,
				SessionLaps: sessions[i].SessionLaps,
				SessionType: sessions[i].SessionType,
				SessionName: sessions[i].SessionName,
			}
		}
	}

	return session
}

func (s *Session) SetState(state int) {
	switch state {
	case invalid:
		s.SessionState = "Invalid"
	case getInCar:
		s.SessionState = "Get In Car"
	case warmup:
		s.SessionState = "Warmup"
	case paradeLaps:
		s.SessionState = "Parade Laps"
	case racing:
		s.SessionState = "Racing"
	case checkered:
		s.SessionState = "Checkered"
	case coolDown:
		s.SessionState = "Cool Down"
	}
}
