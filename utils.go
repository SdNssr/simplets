package main

import (
	"encoding/binary"
	"math"
)

func btof(bytes []byte) float64 {
    bits := binary.LittleEndian.Uint64(bytes)
    return math.Float64frombits(bits)
}

func ftob(float float64) []byte {
    bits := math.Float64bits(float)
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, bits)
    return bytes
}
