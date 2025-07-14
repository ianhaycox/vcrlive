package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLivePositions(t *testing.T) {
	l := LivePositions{Weekend: Weekend{TrackID: 1}}
	expected := "{\n  \"weekend\": {\n    \"track_id\": 1\n  },\n  \"session\": {}\n}"
	assert.Equal(t, expected, l.String())
}
