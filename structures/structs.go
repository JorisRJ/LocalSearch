package structures

import (
	"LocalSearch/utils"
	"image"
	"math/rand"
)

// Pixel is a simple struct holding the info of a pixel
type Pixel struct {
	R, G, B, A uint8
}

//SetBlack sets the pixel to black
func (p *Pixel) SetBlack() {
	p.R = 0
	p.G = 0
	p.B = 0
	p.A = 255
}

// SetPixel sets a pixel to the color of the other
func (p *Pixel) SetPixel(other Pixel) {
	p.R = other.R
	p.G = other.G
	p.B = other.B
	p.A = other.A
}

// A Picture is a picture composted of pixels
type Picture struct {
	Pixels []Pixel
	Width  int
	Height int
}

// BlackPicture gives a new black picture
func BlackPicture(width int, height int) Picture {
	pixels := make([]Pixel, width*height)
	for i := 0; i < width*height; i++ {
		pixels[i] = Pixel{R: 0, G: 0, B: 0, A: 255}
	}

	return Picture{Pixels: pixels, Width: width, Height: height}
}

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
