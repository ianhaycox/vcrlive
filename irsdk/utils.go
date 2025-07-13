package irsdk

import (
	"encoding/binary"
	"math"
	"strings"
)

func byte4ToInt(in []byte) int {
	return int(binary.LittleEndian.Uint32(in))
}

func byte4ToFloat(in []byte) float32 {
	bits := binary.LittleEndian.Uint32(in)
	return math.Float32frombits(bits)
}

func byte8ToFloat(in []byte) float64 {
	bits := binary.LittleEndian.Uint64(in)
	return math.Float64frombits(bits)
}

func bytesToString(in []byte) string {
	return strings.TrimRight(string(in), "\x00")
}
