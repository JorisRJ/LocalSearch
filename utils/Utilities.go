package utils

// MoveMax is the maximum amount a point may be moved by a small mutation
const MoveMax int = 5

// ColorMax is the maximum amount a r,g,b value of a color may be changed by a small mutation
const ColorMax int = 8

// RadiusMax is the max amount the squared radius may be adjusted by a small mutation
const RadiusMax int = 4

// RadiusMaxReroll is the maximum radius that may be chosen in a complete reroll
const RadiusMaxReroll int = 40

// StartCircleAmount is the starting amount of circles in a scene
const StartCircleAmount int = 256

// Clamp is not in the math package sadface
func Clamp(min int, max int, value int) int {
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
