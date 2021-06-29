package structures

import (
	"LocalSearch/utils"
	"fmt"
	"image"
	"math/rand"
	"strconv"
	"strings"
)

/*
Small explainer
If there is a grid of 9x9 squares with triangles, there are 10x10 anchor points,
thus 9x9x2 = 162 triangles and 100 anchor points.
Each of the corners of a triangle (Q1, Q2, Q3) will point to a place in the anchor point array
Thus most anchors are connected to 6 triangles.

*/
// TriangleScene holds anchors and triangles
type TriangleScene struct {
	Anchors           []Coord
	Triangles         []Triangle
	Width, Height     int
	TrWidth, TrHeight int
}

// Coord is a coordinate of float32
type Coord struct {
	X float32
	Y float32
}

// Triangle holds indices to the anchorpoints array and has a color
type Triangle struct {
	Q1, Q2, Q3 int
	Color      Pixel
}

// MutationData holds relevant data returned by the mutation
type MutationData struct {
	Bounds image.Rectangle
	Undo   string
}

// CalculateOuterBounds gives the outer bounds of the triangles that are affected by one anchor
func CalculateOuterBounds(anchor int, triangleScene *TriangleScene) image.Rectangle {
	var coords []Coord

	coords = append(coords, triangleScene.Anchors[anchor])
	coords = append(coords, triangleScene.Anchors[anchor-1])
	coords = append(coords, triangleScene.Anchors[anchor+1])
	coords = append(coords, triangleScene.Anchors[anchor-triangleScene.TrWidth])
	coords = append(coords, triangleScene.Anchors[anchor-triangleScene.TrWidth+1])
	coords = append(coords, triangleScene.Anchors[anchor+triangleScene.TrWidth])
	coords = append(coords, triangleScene.Anchors[anchor+triangleScene.TrWidth-1])

	// Note the reverse
	minX := float32(triangleScene.Width)
	maxX := float32(0)
	minY := float32(triangleScene.Height)
	maxY := float32(0)

	for i := range coords {
		if coords[i].X < minX {
			minX = coords[i].X
		}
		if coords[i].X > maxX {
			minX = coords[i].X
		}
		if coords[i].Y < minY {
			minX = coords[i].Y
		}
		if coords[i].Y < maxY {
			minX = coords[i].Y
		}
	}

	return image.Rect(int(minX), int(minY), int(maxX), int(maxY))
}

// SingleTriangleBounds calculates the bounds around a single triangle
func SingleTriangleBounds(tri Triangle, triangleScene *TriangleScene) image.Rectangle {
	var coords []Coord

	coords = append(coords, triangleScene.Anchors[tri.Q1])
	coords = append(coords, triangleScene.Anchors[tri.Q2])
	coords = append(coords, triangleScene.Anchors[tri.Q3])

	// Note the reverse
	minX := float32(triangleScene.Width)
	maxX := float32(0)
	minY := float32(triangleScene.Height)
	maxY := float32(0)

	for i := range coords {
		if coords[i].X < minX {
			minX = coords[i].X
		}
		if coords[i].X > maxX {
			minX = coords[i].X
		}
		if coords[i].Y < minY {
			minX = coords[i].Y
		}
		if coords[i].Y < maxY {
			minX = coords[i].Y
		}
	}

	return image.Rect(int(minX), int(minY), int(maxX+1), int(maxY+1))
}

// NewTriangleSceneHeadstart creates a gridded triangle picture with poorly guessed colors
func NewTriangleSceneHeadstart(trWidth int, trHeight int, pic *Picture) TriangleScene {
	ySpacing := float32(pic.Height) / float32(trHeight-1)
	xSpacing := float32(pic.Width) / float32(trWidth-1)

	anchors := make([]Coord, trHeight*trWidth)
	for y := 0; y < trHeight; y++ {
		for x := 0; x < trWidth; x++ {
			anchors[x+y*trWidth] = Coord{
				X: float32(x) * xSpacing,
				Y: float32(y) * ySpacing,
			}
		}
	}

	triangles := make([]Triangle, trHeight*trWidth*2)
	for y := 0; y < trHeight-1; y++ {
		for x := 0; x < trWidth-1; x++ {
			triangles[(x+y*trWidth)*2] = Triangle{
				Q1:    y*trWidth + x,
				Q2:    y*trWidth + x + 1,
				Q3:    (y+1)*trWidth + x,
				Color: Pixel{R: 0, G: 0, B: 0, A: 255},
			}
			triangles[(x+y*trWidth)*2+1] = Triangle{
				Q1:    y*trWidth + x + 1,
				Q2:    (y+1)*trWidth + x,
				Q3:    (y+1)*trWidth + x + 1,
				Color: Pixel{R: 0, G: 100, B: 0, A: 255},
			}
		}
	}

	for i := range triangles {
		tr := &triangles[i]
		col := FourPointAverage(anchors[tr.Q1], anchors[tr.Q2], anchors[tr.Q3], *pic)
		tr.Color.SetPixel(col)
	}

	return TriangleScene{
		Anchors:   anchors,
		Triangles: triangles,
		TrWidth:   trWidth,
		TrHeight:  trHeight,
		Width:     pic.Width,
		Height:    pic.Height,
	}
}

// Draw walks accross the sides of each triangle and horizontally draws all the pixels in between
func (trs *TriangleScene) Draw(pic *Picture) {
	maxIndex := len(pic.Pixels)
	var top, mid, bot *Coord
	var q1, q2, q3 *Coord

	// It might hurt to read this
	// This draws each triangle
	for _, tr := range trs.Triangles {
		q1 = &trs.Anchors[tr.Q1]
		q2 = &trs.Anchors[tr.Q2]
		q3 = &trs.Anchors[tr.Q3]

		if q1.Y < q2.Y {
			if q1.Y < q3.Y {
				top = q1

				if q2.Y < q3.Y {
					mid = q2
					bot = q3
				} else {
					mid = q3
					bot = q2
				}
			} else {
				top = q3
				mid = q1
				bot = q2
			}
		} else if q2.Y < q3.Y {
			top = q2
			if q1.Y < q3.Y {
				mid = q1
				bot = q3
			} else {
				mid = q3
				bot = q1
			}
		} else {
			top = q3
			mid = q2
			bot = q1
		}

		// Deltas
		dx1 := mid.X - top.X
		dx2 := bot.X - top.X
		dx3 := bot.X - mid.X
		dy1 := mid.Y - top.Y
		dy2 := bot.Y - top.Y
		dy3 := bot.Y - mid.Y

		// Catch the divide by zero before it happens
		if dy1 == 0 {
			dy1 = 1
		}

		if dy2 == 0 {
			dy2 = 1
		}

		if dy3 == 0 {
			dy3 = 1
		}

		// Slopes
		dxy1 := dx1 / dy1
		dxy2 := dx2 / dy2
		dxy3 := dx3 / dy3

		var x1, x2 float32
		// Determine start point
		if top.Y == mid.Y {
			x1 = utils.FMin(top.X, mid.X)
			x2 = utils.FMax(top.X, mid.X)
		} else {
			x1 = top.X
			x2 = top.X
		}

		y := top.Y

		// Draw the upper part of the triangle
		for y < mid.Y {
			for x := utils.FMin(x1, x2); x < utils.FMax(x1, x2); x++ {

				// This is a failsafe for triangles with slopes that are near horizontal or vertical
				i := int(x) + int(y)*pic.Width
				if i >= maxIndex || i < 0 {
					continue
				}
				pic.Pixels[int(x)+int(y)*pic.Width] = tr.Color
			}

			if y < mid.Y {
				x1 += dxy1
				x2 += dxy2
			}

			y++
		}

		// Draw the lower part of the triangle
		for y <= bot.Y {
			for x := utils.FMin(x1, x2); x < utils.FMax(x1, x2); x++ {

				// This is a failsafe for triangles with slopes that are near horizontal or vertical
				i := int(x) + int(y)*pic.Width
				if i >= maxIndex || i < 0 {
					continue
				}
				// For bugfixing
				// if int(x)+int(y)*pic.Width == 732202 {
				// 	fmt.Printf("\nx: %v; y: %v\n", x, y)
				// 	fmt.Printf("bot.X: %v; bot.Y: %v\n", bot.X, bot.Y)
				// 	fmt.Printf("x1: %v; x2: %v\n", x1, x2)
				// 	fmt.Printf("Width %v\n", pic.Width)
				// 	fmt.Printf("Q1: %v; Q2: %v; Q3: %v\n", q1, q2, q3)
				// }
				pic.Pixels[int(x)+int(y)*pic.Width] = tr.Color
			}

			if y < bot.Y {
				x1 += dxy3
				x2 += dxy2
			}

			y++
		}

	}

}

// Mutate blablabla
func (trs *TriangleScene) Mutate() MutationData {

	r := rand.Intn(2)
	var bounds image.Rectangle
	var undo string

	switch r {
	case 0:
		// Displace by a little bit
		var anc int

		for {
			anc = rand.Intn(len(trs.Anchors))

			// Check to look that it is not a border
			if (anc%trs.TrWidth != 0) && (anc%trs.TrWidth != trs.TrWidth-1) {
				if (anc > trs.TrWidth) && (anc < len(trs.Anchors)-trs.TrWidth) {
					break
				}
			}
		}

		disp := Coord{
			X: rand.Float32()*float32(utils.AnchorMoveMax) - float32(utils.AnchorMoveMax)/2,
			Y: rand.Float32()*float32(utils.AnchorMoveMax) - float32(utils.AnchorMoveMax)/2,
		}

		undo = fmt.Sprintf("disp;%v;%v;%v", anc, trs.Anchors[anc].X, trs.Anchors[anc].Y)

		trs.Anchors[anc].X += disp.X
		trs.Anchors[anc].Y += disp.Y

		trs.Anchors[anc].X = utils.FClamp(0, float32(trs.Width-1), trs.Anchors[anc].X)
		trs.Anchors[anc].Y = utils.FClamp(0, float32(trs.Height-1), trs.Anchors[anc].Y)

		bounds = CalculateOuterBounds(anc, trs)
		break
	case 1:
		// Choose a different collor
		tri := rand.Intn(len(trs.Triangles))

		undo = fmt.Sprintf("recol;%v;%v;%v;%v;%v", tri, trs.Triangles[tri].Color.R, trs.Triangles[tri].Color.G,
			trs.Triangles[tri].Color.B, trs.Triangles[tri].Color.A)

		trs.Triangles[tri].Color.R = utils.Clamp8(0, 255, trs.Triangles[tri].Color.R-uint8(utils.ColorMax)+uint8(rand.Intn(utils.ColorMax/2)))
		trs.Triangles[tri].Color.G = utils.Clamp8(0, 255, trs.Triangles[tri].Color.G-uint8(utils.ColorMax)+uint8(rand.Intn(utils.ColorMax/2)))
		trs.Triangles[tri].Color.B = utils.Clamp8(0, 255, trs.Triangles[tri].Color.B-uint8(utils.ColorMax)+uint8(rand.Intn(utils.ColorMax/2)))
		//trs.Triangles[tri].Color.A = utils.Clamp8(0, 255, trs.Triangles[tri].Color.A - uint8(utils.ColorMax) + uint8(rand.Intn(utils.ColorMax / 2)))
		bounds = SingleTriangleBounds(trs.Triangles[tri], trs)
		break
	case 2:
		break
	case 3:
		break
	case 4:
		break

	default:
		fmt.Println("Yeet")
	}

	return MutationData{
		Bounds: bounds,
		Undo:   undo,
	}
}

// UndoMutation makes a mutation undone
func UndoMutation(undostr string, trs *TriangleScene) {
	s := strings.Split(undostr, ";")

	switch s[0] {
	case "disp":
		anc, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			fmt.Println("Parse failed")
			panic(err)
		}

		x, err := strconv.ParseFloat(s[2], 32)
		if err != nil {
			fmt.Println("X parse failed")
			panic(err)
		}
		trs.Anchors[anc].X = float32(x)

		y, err := strconv.ParseFloat(s[3], 32)
		if err != nil {
			fmt.Println("Y parse failed")
			panic(err)
		}
		trs.Anchors[anc].Y = float32(y)
		break
	case "recol":
		tri, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			fmt.Println("Parse failed")
			panic(err)
		}

		r, err := strconv.ParseFloat(s[2], 32)
		if err != nil {
			fmt.Println("X parse failed")
			panic(err)
		}
		trs.Triangles[tri].Color.R = uint8(r)

		g, err := strconv.ParseFloat(s[3], 32)
		if err != nil {
			fmt.Println("X parse failed")
			panic(err)
		}
		trs.Triangles[tri].Color.G = uint8(g)

		b, err := strconv.ParseFloat(s[4], 32)
		if err != nil {
			fmt.Println("X parse failed")
			panic(err)
		}
		trs.Triangles[tri].Color.B = uint8(b)

		a, err := strconv.ParseFloat(s[5], 32)
		if err != nil {
			fmt.Println("X parse failed")
			panic(err)
		}
		trs.Triangles[tri].Color.A = uint8(a)

		break
	}
}
