package search

import (
	st "LocalSearch/structures"
	"image"
	"math"
)

// Scene is a potato
type Scene interface {
	Mutate()
	Draw(pic *st.Picture)
}

// HELLA SLOW DONT USE
func logFitness(A *st.Picture, B *st.Picture) uint64 {
	var sum uint64

	for i, p := range A.Pixels {
		var partSum float64
		partSum += math.Log(float64((p.R-B.Pixels[i].R)*(p.R-B.Pixels[i].R) + 1))
		partSum += math.Log(float64((p.G-B.Pixels[i].G)*(p.G-B.Pixels[i].G) + 1))
		partSum += math.Log(float64((p.B-B.Pixels[i].B)*(p.B-B.Pixels[i].B) + 1))
		partSum += math.Log(float64((p.A-B.Pixels[i].A)*(p.A-B.Pixels[i].A) + 1))
		sum += uint64(partSum)
	}

	return sum
}

//Big sidenote: this probably is screwed up due to overflow lmao
func squaredFitness(A *st.Picture, B *st.Picture) uint64 {
	var sum uint64

	for i, p := range A.Pixels {
		var partSum uint64
		partSum += uint64((p.R - B.Pixels[i].R) * (p.R - B.Pixels[i].R))
		partSum += uint64((p.G - B.Pixels[i].G) * (p.G - B.Pixels[i].G))
		partSum += uint64((p.B - B.Pixels[i].B) * (p.B - B.Pixels[i].B))
		partSum += uint64((p.A - B.Pixels[i].A) * (p.A - B.Pixels[i].A))
		sum += partSum
	}

	return sum
}

func absFitness(A *st.Picture, B *st.Picture) uint64 {
	var sum uint64

	for i, p := range A.Pixels {
		var partSum uint64
		partSum += uint64(colDiff(p.R, B.Pixels[i].R))
		partSum += uint64(colDiff(p.G, B.Pixels[i].G))
		partSum += uint64(colDiff(p.B, B.Pixels[i].B))
		partSum += uint64(colDiff(p.A, B.Pixels[i].A))
		sum += partSum
	}

	return sum
}

// Big sidenote: this probably is screwed up due to overflow lmao
// Same as above but only for the sections that changed
func absFitnessSubSections(A *st.Picture, B *st.Picture, subSections *[]image.Rectangle) uint64 {
	var sum uint64

	for _, rect := range *subSections {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				i := y*A.Width + x
				var partSum uint64
				partSum += uint64(colDiff(A.Pixels[i].R, B.Pixels[i].R))
				partSum += uint64(colDiff(A.Pixels[i].G, B.Pixels[i].G))
				partSum += uint64(colDiff(A.Pixels[i].B, B.Pixels[i].B))
				partSum += uint64(colDiff(A.Pixels[i].A, B.Pixels[i].A))
				sum += partSum
			}
		}
	}

	return sum
}

func colDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}
