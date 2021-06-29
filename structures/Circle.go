package structures

import (
	"LocalSearch/utils"
	"image"
	"math/rand"
)

// Circle holds the basic info of a circle,
type Circle struct {
	Radius, X, Y int
	Color        Pixel
}

// CalcBounds calculates the outer most points of two shapes
func CalcBounds(circ Circle, bounds image.Rectangle, picWidth int, picHeight int) image.Rectangle {
	xmin := circ.X - circ.Radius
	xmax := circ.X + circ.Radius
	ymin := circ.Y - circ.Radius
	ymax := circ.Y + circ.Radius

	xmin = utils.Min(xmin, bounds.Min.X)
	xmax = utils.Max(xmax, bounds.Max.X)
	ymin = utils.Min(ymin, bounds.Min.Y)
	ymax = utils.Max(ymax, bounds.Max.Y)

	return image.Rect(
		utils.Clamp(0, picWidth, xmin),
		utils.Clamp(0, picHeight, ymin),
		utils.Clamp(0, picWidth, xmax),
		utils.Clamp(0, picHeight, ymax),
	)
}

// SkipRect is a rectangle used to have no impact on the bounds
func SkipRect() image.Rectangle {
	return image.Rect(2000000, 2000000, -2000000, -2000000)
}

// NewRandomCircle Creates a new random circle
func NewRandomCircle(width int, height int) Circle {
	c := Circle{
		Radius: rand.Intn(utils.RadiusMaxReroll),
		X:      rand.Intn(width),
		Y:      rand.Intn(height),
		Color: Pixel{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			//A: uint8(rand.Intn(255)),
			A: 255,
		},
	}

	return c
}
