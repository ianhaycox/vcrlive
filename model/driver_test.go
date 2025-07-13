package model

import (
	"testing"

	"github.com/ianhaycox/vcrlive/irsdk/iryaml"
	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	t.Run("New drivers should return map", func(t *testing.T) {
		irDrivers := []iryaml.Driver{
			{CarIdx: 0, UserName: "Pace car"},
			{CarIdx: 1, UserName: "Driver 1"},
			{CarIdx: 2, UserName: "Driver 2"},
			{CarIdx: 3, UserName: "Driver 3"},
			{CarIdx: 4, UserName: "Driver 4"},
			{CarIdx: 0, UserName: ""},
		}

		d := NewDrivers(irDrivers)
		d.SetPositions([]int{0, 2, 6, 13, 1, 0, 0, 0, 0})
		d.SetLaps([]int{0, 5, 6, 7, 8, 0, 0, 0, 0})

		expectedDrivers := Drivers{
			1: Driver{CarIdx: 1, UserName: "Driver 1", ClassPosition: 2, Lap: 5},
			2: Driver{CarIdx: 2, UserName: "Driver 2", ClassPosition: 6, Lap: 6},
			3: Driver{CarIdx: 3, UserName: "Driver 3", ClassPosition: 13, Lap: 7},
			4: Driver{CarIdx: 4, UserName: "Driver 4", ClassPosition: 1, Lap: 8},
		}

		assert.Equal(t, expectedDrivers, d)
	})
}
