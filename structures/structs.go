package structures

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

// A Picture is a picture composed of pixels
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

// FourPointAverage gives a color averaged between the midmid point and the halfwaypoints between
// the middle and of the line from midmid to each corner
func FourPointAverage(q1, q2, q3 Coord, pic Picture) Pixel {
	midx := (q1.X + q2.X + q3.X) / 3
	midy := (q1.Y + q2.Y + q3.Y) / 3

	// Halfway point between corner (q) and the middle
	hq1x := (q1.X + midx) / 2
	hq1y := (q1.Y + midy) / 2
	hq2x := (q2.X + midx) / 2
	hq2y := (q2.Y + midy) / 2
	hq3x := (q3.X + midx) / 2
	hq3y := (q3.Y + midy) / 2

	midcol := pic.Pixels[int(midx)+int(midy)*pic.Width]
	hq1col := pic.Pixels[int(hq1x)+int(hq1y)*pic.Width]
	hq2col := pic.Pixels[int(hq2x)+int(hq2y)*pic.Width]
	hq3col := pic.Pixels[int(hq3x)+int(hq3y)*pic.Width]

	return Pixel{
		R: (midcol.R >> 2) + (hq1col.R >> 2) + (hq2col.R >> 2) + (hq3col.R >> 2),
		G: (midcol.G >> 2) + (hq1col.G >> 2) + (hq2col.G >> 2) + (hq3col.G >> 2),
		B: (midcol.B >> 2) + (hq1col.B >> 2) + (hq2col.B >> 2) + (hq3col.B >> 2),
		A: (midcol.A >> 2) + (hq1col.A >> 2) + (hq2col.A >> 2) + (hq3col.A >> 2),
	}
}
