package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLivePositions(t *testing.T) {
	l := LivePositions{Weekend: Weekend{TrackID: 1}}
	expected := "{\n  \"weekend\": {\n    \"track_id\": 1,\n    \"track_display_name\": \"\",\n    \"track_config_name\": \"\",\n    \"series_id\": 0,\n    \"season_id\": 0,\n    \"session_id\": 0,\n    \"sub_session_id\": 0,\n    \"official\": 0,\n    \"race_week\": 0,\n    \"event_type\": \"\",\n    \"category\": \"\",\n    \"num_car_classes\": 0,\n    \"num_car_types\": 0\n  },\n  \"session\": {\n    \"session_num\": 0,\n    \"session_laps\": \"\",\n    \"session_type\": \"\",\n    \"session_name\": \"\",\n    \"session_state\": \"\",\n    \"error_text\": \"\"\n  }\n}"
	assert.Equal(t, expected, l.String())
}
