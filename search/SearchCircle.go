package search

import (
	st "LocalSearch/structures"
	"LocalSearch/utils"
	"fmt"
	"image"
	"math/rand"
)

// CircScene holds all the circles to implement the
type CircScene struct {
	Circles []st.Circle
	Width   int
	Height  int
}

// NewRandomScene creates a new random scene
func NewRandomScene(width int, height int) CircScene {
	circles := make([]st.Circle, utils.StartCircleAmount)

	for i := 0; i < utils.StartCircleAmount; i++ {
		circles[i] = st.NewRandomCircle(width, height)
	}

	return CircScene{Circles: circles, Width: width, Height: height}
}

// HeadStartedScene creates a scene with colors picked from the original
func HeadStartedScene(width int, height int, original *st.Picture) CircScene {
	circles := make([]st.Circle, utils.StartCircleAmount)

	for i := 0; i < utils.StartCircleAmount; i++ {
		circles[i] = st.NewRandomCircle(width, height)
		circles[i].Color = original.Pixels[circles[i].Y*width+circles[i].X]
	}

	return CircScene{Circles: circles, Width: width, Height: height}
}

// Mutate alters the scene by one mutation, returns the subsection of the image that changed
func (circScene *CircScene) Mutate() []image.Rectangle {
	width := circScene.Width
	height := circScene.Height
	r := rand.Intn(6)
	switch r {
	case 0:
		// Small movement change
		circ := rand.Intn(len(circScene.Circles))
		bounds := st.CalcBounds(circScene.Circles[circ], st.SkipRect(), width, height)
		circScene.Circles[circ].X = circScene.Circles[circ].X + rand.Intn(2*utils.MoveMax) - utils.MoveMax
		circScene.Circles[circ].Y = circScene.Circles[circ].Y + rand.Intn(2*utils.MoveMax) - utils.MoveMax
		bounds = st.CalcBounds(circScene.Circles[circ], bounds, width, height)
		return []image.Rectangle{bounds}

	case 1:
		// Complete relocation
		circ := rand.Intn(len(circScene.Circles))
		boundsOld := st.CalcBounds(circScene.Circles[circ], st.SkipRect(), width, height)
		circScene.Circles[circ].X = rand.Intn(circScene.Width)
		circScene.Circles[circ].Y = rand.Intn(circScene.Height)
		boundsNew := st.CalcBounds(circScene.Circles[circ], st.SkipRect(), width, height)

		return []image.Rectangle{boundsOld, boundsNew}

	case 2:
		// Small color change
		circ := rand.Intn(len(circScene.Circles))
		bounds := st.CalcBounds(circScene.Circles[circ], st.SkipRect(), width, height)
		circScene.Circles[circ].Color.R = uint8(utils.Clamp(0, 255, int(circScene.Circles[circ].Color.R)+rand.Intn(2*utils.ColorMax)-int(utils.ColorMax)))
		circScene.Circles[circ].Color.G = uint8(utils.Clamp(0, 255, int(circScene.Circles[circ].Color.G)+rand.Intn(2*utils.ColorMax)-int(utils.ColorMax)))
		circScene.Circles[circ].Color.B = uint8(utils.Clamp(0, 255, int(circScene.Circles[circ].Color.B)+rand.Intn(2*utils.ColorMax)-int(utils.ColorMax)))
		//circScene.Circles[circ].Color.A = uint8(utils.Clamp(0, 255, int(circScene.Circles[circ].Color.A)+rand.Intn(2*utils.ColorMax)-int(utils.ColorMax)))
		return []image.Rectangle{bounds}

	case 3:
		// Reroll color
		circ := rand.Intn(len(circScene.Circles))
		bounds := st.CalcBounds(circScene.Circles[circ], st.SkipRect(), width, height)
		circScene.Circles[circ].Color.R = uint8(rand.Intn(255))
		circScene.Circles[circ].Color.G = uint8(rand.Intn(255))
		circScene.Circles[circ].Color.B = uint8(rand.Intn(255))
		//circScene.Circles[circ].Color.A = uint8(rand.Intn(255))
		return []image.Rectangle{bounds}

	case 4:
		// Small radius change
		circ := rand.Intn(len(circScene.Circles))
		bounds := st.CalcBounds(circScene.Circles[circ], st.SkipRect(), width, height)
		circScene.Circles[circ].Radius = circScene.Circles[circ].Radius + rand.Intn(2*utils.RadiusMax) - utils.RadiusMax
		bounds = st.CalcBounds(circScene.Circles[circ], bounds, width, height)
		return []image.Rectangle{bounds}

	case 5:
		// Complete radius reroll
		circ := rand.Intn(len(circScene.Circles))
		bounds := st.CalcBounds(circScene.Circles[circ], st.SkipRect(), width, height)
		circScene.Circles[circ].Radius = rand.Intn(utils.RadiusMaxReroll)
		bounds = st.CalcBounds(circScene.Circles[circ], bounds, width, height)
		return []image.Rectangle{bounds}

	case 6:
		// Layer Swap small
		break
	case 7:
		// Layer swap Random
		break
	case 8:
		// Add circle
		break
	case 9:
		// Delete circle
		break
	}

	return []image.Rectangle{}
}

// Draw draws the scene to a picture
func (circScene *CircScene) Draw(pic *st.Picture) {
	// Loop over each pixel, within each pixel loop over the circles

	for _, circle := range circScene.Circles {
		xmin := utils.Clamp(0, pic.Width, circle.X-circle.Radius)
		xmax := utils.Clamp(0, pic.Width, circle.X+circle.Radius)
		ymin := utils.Clamp(0, pic.Height, circle.Y-circle.Radius)
		ymax := utils.Clamp(0, pic.Height, circle.Y+circle.Radius)

		for y := ymin; y < ymax; y++ {
			for x := xmin; x < xmax; x++ {
				distToCenter := (circle.X-x)*(circle.X-x) + (circle.Y-y)*(circle.Y-y)
				if distToCenter < (circle.Radius * circle.Radius) {
					pic.Pixels[y*pic.Width+x].SetPixel(circle.Color)
				}
			}
		}
	}
}

// AreaDraw draws only a selected area HELLA SLOW
func (circScene *CircScene) AreaDraw(pic *st.Picture, subSections []image.Rectangle) {

	for _, rect := range subSections {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				// Go through in reverse due to layered draw
				for i := len(circScene.Circles) - 1; i >= 0; i-- {
					circle := circScene.Circles[i]
					distToCenter := (circle.X-x)*(circle.X-x) + (circle.Y-y)*(circle.Y-y)
					if distToCenter < (circle.Radius * circle.Radius) {
						pic.Pixels[y*pic.Width+x].SetPixel(circle.Color)
						break
					}
				}
			}
		}
	}
}

// MutationRounds mutates for a number of rounds
func MutationRounds(rounds int, original *st.Picture, match *st.Picture, scene *CircScene) {
	for i := 0; i < rounds; i++ {
		if i%10 == 0 {
			fmt.Printf("\rProgress: %v%%", (i+1)*100/rounds)
		}
		//fmt.Print(".")

		var dup CircScene
		utils.Clone(scene, &dup)
		fitOld := AbsFitness(original, match)
		scene.Mutate()
		scene.Draw(match)
		fitNew := AbsFitness(original, match)

		if fitNew < fitOld {
			continue
		} else {
			scene = &dup
			scene.Draw(match)
		}
	}
}
