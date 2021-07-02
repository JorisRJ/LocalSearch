package search

import (
	st "LocalSearch/structures"
	"fmt"
	"image"
)

// MutationRoundsTriangles mutates for a number of rounds
func MutationRoundsTriangles(rounds int, original *st.Picture, match *st.Picture, scene *st.TriangleScene) {
	for i := 0; i < rounds; i++ {
		if i%10 == 0 {
			fmt.Printf("\rProgress: %v%%", (i+1)*100/rounds)
		}

		muData := scene.Mutate()
		fitOld := absFitnessSubSections(original, match, &[]image.Rectangle{muData.Bounds})
		scene.Draw(match)
		fitNew := absFitnessSubSections(original, match, &[]image.Rectangle{muData.Bounds})

		if fitNew < fitOld {
			continue
		} else {
			st.UndoMutation(muData.Undo, scene)
			scene.Draw(match)
		}
	}
}

// Alright algoritme idee
// Goroutines very many
