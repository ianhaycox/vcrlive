package model

import "github.com/ianhaycox/vcrlive/irsdk/iryaml"

type Weekend struct {
	TrackID          int    `json:"track_id"`
	TrackDisplayName string `json:"track_display_name"`
	TrackConfigName  string `json:"track_config_name"`
	SeriesID         int    `json:"series_id"`
	SeasonID         int    `json:"season_id"`
	SessionID        int    `json:"session_id"`
	SubSessionID     int    `json:"sub_session_id"`
	Official         int    `json:"official"`
	RaceWeek         int    `json:"race_week"`
	EventType        string `json:"event_type"`
	Category         string `json:"category"`
	NumCarClasses    int    `json:"num_car_classes"`
	NumCarTypes      int    `json:"num_car_types"`
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
