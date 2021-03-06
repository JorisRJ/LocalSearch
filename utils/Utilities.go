package utils

import (
	"bytes"
	"encoding/gob"
)

// AnchorMoveMax is the maximum amount a circle may be moved by a small mutation
const AnchorMoveMax int = 5

// ColorMax is the maximum amount a r,g,b value of a color may be changed by a small mutation
const ColorMax int = 8

// Clamp is not in the math package sadface (or I did not look hard enough)
func Clamp(min int, max int, value int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}

	return value
}

// Clamp8 is clamp for uint8
func Clamp8(min uint8, max uint8, value uint8) uint8 {
	if value < min {
		return min
	} else if value > max {
		return max
	}

	return value
}

// Min is not there for ints
func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max is not there for ints
func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clone clones; yes the floor is made out of floor
func Clone(a, b interface{}) {

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
}

// FMin gives the smallest float
func FMin(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// FMax gives the largest float
func FMax(a, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

// LowestIndex returns the index of the lowest value in the list
func LowestIndex(list []uint64) int {
	low := list[0]
	index := 0
	for i := range list {
		if list[i] < low {
			index = i
		}
	}
	return index
}
