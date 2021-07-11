package search

import (
	st "LocalSearch/structures"
	"image"
)

// AbsFitness is the fitness of a picture using the formula abs(A-B)
func AbsFitness(A *st.Picture, B *st.Picture) uint64 {
	var sum uint64

	for i, p := range A.Pixels {
		var partSum uint64
		partSum += uint64(ColDiff(p.R, B.Pixels[i].R))
		partSum += uint64(ColDiff(p.G, B.Pixels[i].G))
		partSum += uint64(ColDiff(p.B, B.Pixels[i].B))
		partSum += uint64(ColDiff(p.A, B.Pixels[i].A))
		sum += partSum
	}

	return sum
}

// Calculates the fitness of a picture by testing it against another using the absolute() function
func absFitnessSubSections(A *st.Picture, B *st.Picture, subSections *[]image.Rectangle) uint64 {
	var sum uint64

	for _, rect := range *subSections {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				i := y*A.Width + x
				var partSum uint64
				partSum += uint64(ColDiff(A.Pixels[i].R, B.Pixels[i].R))
				partSum += uint64(ColDiff(A.Pixels[i].G, B.Pixels[i].G))
				partSum += uint64(ColDiff(A.Pixels[i].B, B.Pixels[i].B))
				partSum += uint64(ColDiff(A.Pixels[i].A, B.Pixels[i].A))
				sum += partSum
			}
		}
	}

	return sum
}

// ColDiff gives the difference in color
func ColDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a

}
