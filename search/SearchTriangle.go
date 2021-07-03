package search

import (
	st "LocalSearch/structures"
	"LocalSearch/utils"
	"fmt"
	"image"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// MutationRoundsTriangles mutates for a number of rounds
func MutationRoundsTriangles(rounds int, original *st.Picture, match *st.Picture, scene *st.TriangleScene) {
	random := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < rounds; i++ {
		if i%10 == 0 {
			fmt.Printf("\rProgress: %v%%", (i+1)*100/rounds)
		}

		var fitOld, fitNew uint64
		muData := scene.Mutate(random)

		//If the mutation is a recolor, only one triangle needs to be recolored
		undo := strings.Split(muData.Undo, ";")
		if undo[0] == "recol" {
			tri, _ := strconv.ParseInt(undo[1], 10, 64)
			fitOld = SingleTriangleFitness(original, match, int(tri), scene)
			scene.SingleTriangleDraw(match, int(tri))
			fitNew = SingleTriangleFitness(original, match, int(tri), scene)
		} else {
			fitOld = absFitnessSubSections(original, match, &[]image.Rectangle{muData.Bounds})
			scene.Draw(match)
			fitNew = absFitnessSubSections(original, match, &[]image.Rectangle{muData.Bounds})
		}

		if fitNew < fitOld {
			continue
		} else {
			st.UndoMutation(muData.Undo, scene)
			undo := strings.Split(muData.Undo, ";")

			if undo[0] == "recol" {
				tri, _ := strconv.ParseInt(undo[1], 10, 64)
				scene.SingleTriangleDraw(match, int(tri))
			} else {
				scene.Draw(match)
			}
		}
	}
}

// Alright algoritme idee
// Goroutines very many

// SingleTriangleFitness checks a single triangle, usefull for a color change
func SingleTriangleFitness(A *st.Picture, B *st.Picture, tri int, trs *st.TriangleScene) uint64 {
	var top, mid, bot *st.Coord
	var q1, q2, q3 *st.Coord
	var sum uint64

	// It might hurt to read this
	// This draws each triangle
	tr := trs.Triangles[tri]
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
	dxy1 := float32(dx1) / float32(dy1)
	dxy2 := float32(dx2) / float32(dy2)
	dxy3 := float32(dx3) / float32(dy3)

	var x1, x2 float32
	var tempX1, tempX2 int
	// Determine start point
	if top.Y == mid.Y {
		tempX1 = utils.Min(top.X, mid.X)
		tempX2 = utils.Max(top.X, mid.X)
	} else {
		tempX1 = top.X
		tempX2 = top.X
	}

	x1 = float32(tempX1)
	x2 = float32(tempX2)

	y := float32(top.Y)

	// Draw the upper part of the triangle
	for y < float32(mid.Y) {
		for x := utils.FMin(x1, x2); x < utils.FMax(x1, x2); x++ {

			i := utils.Clamp(0, A.Width*A.Height-1, int(x+y*float32(A.Width)))
			var partSum uint64
			partSum += uint64(ColDiff(A.Pixels[i].R, B.Pixels[i].R))
			partSum += uint64(ColDiff(A.Pixels[i].G, B.Pixels[i].G))
			partSum += uint64(ColDiff(A.Pixels[i].B, B.Pixels[i].B))
			partSum += uint64(ColDiff(A.Pixels[i].A, B.Pixels[i].A))
			sum += partSum
		}

		if y < float32(mid.Y) {
			x1 += dxy1
			x2 += dxy2
		}

		y++
	}

	// Draw the lower part of the triangle
	for y <= float32(bot.Y) {
		for x := utils.FMin(x1, x2); x < utils.FMax(x1, x2); x++ {

			i := utils.Clamp(0, A.Width*A.Height-1, int(x+y*float32(A.Width)))
			var partSum uint64
			partSum += uint64(ColDiff(A.Pixels[i].R, B.Pixels[i].R))
			partSum += uint64(ColDiff(A.Pixels[i].G, B.Pixels[i].G))
			partSum += uint64(ColDiff(A.Pixels[i].B, B.Pixels[i].B))
			partSum += uint64(ColDiff(A.Pixels[i].A, B.Pixels[i].A))
			sum += partSum
		}

		if y < float32(bot.Y) {
			x1 += dxy3
			x2 += dxy2
		}

		y++
	}

	return sum

}
