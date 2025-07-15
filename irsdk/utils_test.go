package irsdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	t.Run("int conversions", func(t *testing.T) {
		assert.Equal(t, 0, byte4ToInt([]byte{0, 0, 0, 0}))
		assert.Equal(t, 4, byte4ToInt([]byte{4, 0, 0, 0}))
		assert.Equal(t, -1, byte4ToInt([]byte{255, 255, 255, 255}))
	})
}
