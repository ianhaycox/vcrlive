package model

import "github.com/ianhaycox/vcrlive/irsdk/iryaml"

const (
	Invalid = iota
	GetInCar
	Warmup
	ParadeLaps
	Racing
	Checkered
	CoolDown
)

type Session struct {
	SessionNum   int    `json:"session_num,omitempty"`
	SessionLaps  string `json:"session_laps,omitempty"`
	SessionType  string `json:"session_type,omitempty"`
	SessionName  string `json:"session_name,omitempty"`
	SessionState string `json:"session_state,omitempty"`
	ErrorText    string `json:"error_text,omitempty"` // Set when SessionState is invalid so the consumer knows about a problem.
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
	case Invalid:
		s.SessionState = "Invalid"
	case GetInCar:
		s.SessionState = "Get In Car"
	case Warmup:
		s.SessionState = "Warmup"
	case ParadeLaps:
		s.SessionState = "Parade Laps"
	case Racing:
		s.SessionState = "Racing"
	case Checkered:
		s.SessionState = "Checkered"
	case CoolDown:
		s.SessionState = "Cool Down"
	}
}
