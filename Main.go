package main

import (
	"LocalSearch/imagerelated"
	"LocalSearch/search"
	"LocalSearch/structures"
	st "LocalSearch/structures"
	"LocalSearch/utils"
	"fmt"
	"image"
	"sync"
	"time"
)

// Threads is the amount of threads the program starts at
const Threads int = 8

// RoundsPerWave is the amount of rounds per wave, every wave some threads will be killed and some duplicated
const RoundsPerWave int = 500000

// AnchorsWidth is the amound of Anchors horizontally
const AnchorsWidth int = 100

// AnchorsHeight is the amount of Anchors vertically
const AnchorsHeight int = 75

// var cpuprofile = flag.String("cpuprofile", "profile", "write cpu profile to file")

func main() {

	// // Use this in case you want to further optimise the code (you need pprof and graphviz)
	// flag.Parse()
	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	pprof.StartCPUProfile(f)
	// 	defer pprof.StopCPUProfile()
	// }

	// Choose one of these three and thange the name to your picture
	singleThreaded("pindakaas", nil)
	// multiThreaded("car")
	//parallelPictures([]string{"pindakaas", "nederland"})
}

// singleThreaded spends one thread on one picture
func singleThreaded(name string, wg *sync.WaitGroup) {

	ogImg, err := imagerelated.OpenImage(fmt.Sprintf("pictures/%v.png", name))
	if err != nil {
		fmt.Print(err)
		return
	}

	// Create an empty picture
	width := ogImg.Bounds().Max.X
	height := ogImg.Bounds().Max.Y
	original := st.Picture{
		Pixels: imagerelated.ImageToPixels(ogImg),
		Width:  width,
		Height: height,
	}

	// Create a new scene
	scene := st.NewTriangleSceneHeadstart(AnchorsWidth, AnchorsHeight, &original)

	// // Open a saved scene
	// scene := *structures.LoadTriangleScene(fmt.Sprintf("saves/%v_save.json", name))

	matchImg := st.BlackPicture(width, height)
	scene.Draw(&matchImg)

	// // Uncomment to see the headstarted scene
	// img1 := imagerelated.PixelsToImage(matchImg.Pixels, image.Rect(0, 0, width, height))
	// imagerelated.SaveImage(img1, "pictures/Headstartedscene.png")

	fmt.Println(fmt.Sprintf("Started on %v", name))

	start := time.Now()

	search.MutationRoundsTriangles(RoundsPerWave, &original, &matchImg, &scene, nil)

	// Info about the speed
	elapsed := time.Since(start)
	fmt.Printf("\n%v Rounds took %s on %v", RoundsPerWave, elapsed, name)
	fmt.Printf("\n%.1f FPS\n", float32(RoundsPerWave)/float32(elapsed.Seconds()))

	// Saving the image and the scene
	img2 := imagerelated.PixelsToImage(matchImg.Pixels, image.Rect(0, 0, width, height))
	imagerelated.SaveImage(img2, fmt.Sprintf("pictures/Triangle_%v.png", name))
	structures.SaveTriangleScene(&scene, fmt.Sprintf("saves/%v_save.json", name))

	if wg != nil {
		wg.Done()
	}
}

// mutlithreaded spends multiple threads on one picture, racing them until one is left.
// It is by far a good algorithm, it is merely some code to get started with multithreading.
func multiThreaded(name string) {

	// Initialize the arrays holding the variables
	originalPictures := make([]st.Picture, Threads)
	matchPictures := make([]st.Picture, Threads)
	scenes := make([]st.TriangleScene, Threads)
	skip := make([]bool, Threads)
	fitnesses := make([]uint64, Threads)
	winner := 0

	ogImg, err := imagerelated.OpenImage(fmt.Sprintf("pictures/%v.png", name))
	if err != nil {
		fmt.Print(err)
		return
	}

	// Create an empty picture
	width := ogImg.Bounds().Max.X
	height := ogImg.Bounds().Max.Y
	original := st.Picture{
		Pixels: imagerelated.ImageToPixels(ogImg),
		Width:  width,
		Height: height,
	}

	// Create copies of the original
	for i := 0; i < Threads; i++ {
		utils.Clone(&original, &originalPictures[i])
	}

	// Create a new scene
	// scene := st.NewTriangleSceneHeadstart(AnchorsWidth, AnchorsHeight, &original)

	// Open a saved scene
	scene := *structures.LoadTriangleScene(fmt.Sprintf("saves/%v_save.json", name))

	matchImg := st.BlackPicture(width, height)
	scene.Draw(&matchImg)

	// Create copies of the match images
	for i := 0; i < Threads; i++ {
		utils.Clone(&matchImg, &matchPictures[i])
	}

	// Create copies of the match images
	for i := 0; i < Threads; i++ {
		utils.Clone(&scene, &scenes[i])
	}

	fmt.Println(fmt.Sprintf("Started %v threads on %v", Threads, name))

	start := time.Now()
	var wg sync.WaitGroup

	// The big loop, every wave a thread is killed, thus there are Thread waves
	for i := 0; i < Threads-1; i++ {

		// Start all goroutines
		for j := 0; j < Threads; j++ {
			if skip[j] {
				continue
			}
			wg.Add(1)

			go search.MutationRoundsTriangles(RoundsPerWave, &originalPictures[i], &matchPictures[i], &scenes[i], &wg)
		}

		wg.Wait()

		// Determine the fitness of the remaining threads
		for j := 0; j < Threads; j++ {
			if skip[j] {
				// Set to the maximum value
				a := -1
				fitnesses[j] = uint64(a)
				continue
			}

			fitnesses[j] = search.AbsFitness(&original, &matchPictures[j])
		}

		// Skip the worst picture
		skip[utils.LowestIndex(fitnesses)] = true
		fmt.Printf("Threads alive: %v\n", (Threads - i - 1))

	}

	// Determine the winner of the threads
	for j := 0; j < Threads; j++ {
		if !skip[j] {
			winner = j
			break
		}
	}

	// The multithreaded rounds are span not work
	elapsed := time.Since(start)
	fmt.Printf("\n%v (Multithreaded) Rounds took %s", Threads*RoundsPerWave, elapsed)
	fmt.Printf("\n%.1f FPS\n", float32(Threads*RoundsPerWave)/float32(elapsed.Seconds()))

	img2 := imagerelated.PixelsToImage(matchPictures[winner].Pixels, image.Rect(0, 0, width, height))
	imagerelated.SaveImage(img2, fmt.Sprintf("pictures/Triangle_%v.png", name))
	structures.SaveTriangleScene(&scene, fmt.Sprintf("saves/%v_save.json", name))
}

// parallelPictures initiates a searcher (and thus a thread) per picture.
// Thus allowing 5 pictures to be searched in parallel.
func parallelPictures(names []string) {
	var wg sync.WaitGroup
	wg.Add(len(names))

	for i := range names {
		go singleThreaded(names[i], &wg)
	}

	wg.Wait()
}
