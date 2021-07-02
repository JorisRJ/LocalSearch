package main

import (
	"LocalSearch/imagerelated"
	"LocalSearch/search"
	st "LocalSearch/structures"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "profile", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ogImg, err := imagerelated.OpenImage("pictures/mountain-lake.png")
	if err != nil {
		fmt.Print(err)
		return
	}

	width := ogImg.Bounds().Max.X
	height := ogImg.Bounds().Max.Y
	original := st.Picture{
		Pixels: imagerelated.ImageToPixels(ogImg),
		Width:  width,
		Height: height,
	}

	scene := st.NewTriangleSceneHeadstart(64, 48, &original)
	//scene := *structures.LoadTriangleScene("saves/triangles1.json")
	matchImg := st.BlackPicture(width, height)
	scene.Draw(&matchImg)
	//img1 := imagerelated.PixelsToImage(matchImg.Pixels, image.Rect(0, 0, width, height))
	//imagerelated.SaveImage(img1, "pictures/Test1_Triangle2.png")

	fmt.Println("Started")

	start := time.Now()

	rounds := 1000
	search.MutationRoundsTriangles(rounds, &original, &matchImg, &scene)

	elapsed := time.Since(start)
	fmt.Printf("\n%v Rounds took %s", rounds, elapsed)
	fmt.Printf("\n%.1f FPS\n", float32(rounds)/float32(elapsed.Seconds()))

	img2 := imagerelated.PixelsToImage(matchImg.Pixels, image.Rect(0, 0, width, height))
	imagerelated.SaveImage(img2, "pictures/Test2_Triangle4.png")
	//structures.SaveTriangleScene(&scene, "saves/triangles1.json")

}
