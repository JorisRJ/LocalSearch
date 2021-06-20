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
	// width := 500
	// height := 500
	// p := st.BlackPicture(width, height)
	// scene := search.NewRandomScene(width, height)
	// scene.Draw(&p)

	// img := imagerelated.PixelsToImage(p.Pixels, image.Rect(0, 0, width, height))
	// imagerelated.SaveImage(img, "Test1.jpg")
	ogImg, err := imagerelated.OpenImage("tree.png")
	if err != nil {
		fmt.Print(err)
		return
	}
	width := ogImg.Bounds().Max.X
	height := ogImg.Bounds().Max.Y
	original := st.Picture{Pixels: imagerelated.ImageToPixels(ogImg), Width: width, Height: height}
	scene := search.NewRandomScene(width, height)
	matchImg := st.BlackPicture(width, height)
	scene.Draw(&matchImg)
	img1 := imagerelated.PixelsToImage(matchImg.Pixels, image.Rect(0, 0, width, height))
	imagerelated.SaveImage(img1, "Test1.jpg")
	fmt.Println("Started")

	start := time.Now()

	search.MutationRounds(1000, &original, &matchImg, &scene)

	elapsed := time.Since(start)
	log.Printf("\n1000 Rounds took %s", elapsed)

	img2 := imagerelated.PixelsToImage(matchImg.Pixels, image.Rect(0, 0, width, height))
	imagerelated.SaveImage(img2, "Test2.jpg")

	// for i := 0; i < 100; i++ {
	// 	fmt.Print(".")
	// 	scene.Mutate()
	// }
	// scene.Draw(&p)
	// img = imagerelated.PixelsToImage(p.Pixels, image.Rect(0, 0, 500, 500))
	// imagerelated.SaveImage(img, "Test2.jpg")
}
