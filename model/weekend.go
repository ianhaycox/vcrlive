package model

import "github.com/ianhaycox/vcrlive/irsdk/iryaml"

type Weekend struct {
	TrackID          int    `json:"track_id,omitempty"`
	TrackDisplayName string `json:"track_display_name,omitempty"`
	TrackConfigName  string `json:"track_config_name,omitempty"`
	SeriesID         int    `json:"series_id,omitempty"`
	SeasonID         int    `json:"season_id,omitempty"`
	SessionID        int    `json:"session_id,omitempty"`
	SubSessionID     int    `json:"sub_session_id,omitempty"`
	Official         int    `json:"official,omitempty"`
	RaceWeek         int    `json:"race_week,omitempty"`
	EventType        string `json:"event_type,omitempty"`
	Category         string `json:"category,omitempty"`
	NumCarClasses    int    `json:"num_car_classes,omitempty"`
	NumCarTypes      int    `json:"num_car_types,omitempty"`
}

func NewWeekend(weekend *iryaml.WeekendInfo) Weekend {
	return Weekend{
		TrackID:          weekend.TrackID,
		TrackDisplayName: weekend.TrackDisplayName,
		TrackConfigName:  weekend.TrackConfigName,
		SeriesID:         weekend.SeriesID,
		SeasonID:         weekend.SeasonID,
		SubSessionID:     weekend.SubSessionID,
		Official:         weekend.Official,
		RaceWeek:         weekend.RaceWeek,
		EventType:        weekend.EventType,
		Category:         weekend.Category,
		NumCarClasses:    weekend.NumCarClasses,
		NumCarTypes:      weekend.NumCarTypes,
	}
}
