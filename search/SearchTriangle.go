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
		//fmt.Print(".")

		//var dup st.TriangleScene
		//utils.Clone(scene, &dup)
		// fitOld := absFitness(original, match)
		muData := scene.Mutate()
		fitOld := absFitnessSubSections(original, match, &[]image.Rectangle{muData.Bounds})
		scene.Draw(match)
		fitNew := absFitnessSubSections(original, match, &[]image.Rectangle{muData.Bounds})
		// fitNew := absFitness(original, match)

		if fitNew < fitOld {
			//fmt.Printf("\nBETTER: %v    %v", fitOld, fitNew)
			//fmt.Println(fitOld - fitNew)
			continue
		} else {
			//fmt.Printf("\n WORSE: %v    %v", fitOld, fitNew)
			//fmt.Println(fitNew - fitOld)
			//scene = &dup
			st.UndoMutation(muData.Undo, scene)
			scene.Draw(match)
		}
	}
}
