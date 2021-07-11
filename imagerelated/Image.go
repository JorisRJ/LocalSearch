package imagerelated

import (
	"LocalSearch/structures"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
)

// OpenImage opens an image file
func OpenImage(path string) (*image.NRGBA, error) {
	fimg, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not open img")
	}

	defer fimg.Close()

	// Decode the image
	img, _, err := image.Decode(fimg)

	if err != nil {

	}

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not Decode img")
	}

	// Works for 32 bit PNG
	if rgbaImg, ok := img.(*image.NRGBA); ok {
		return rgbaImg, nil
	}

	// Do it the slow way
	// This is for 24 bit PNG and maybe jpeg, idk
	nImg := image.NewNRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			offset := nImg.Stride*y + x*4
			r, g, b, a := img.At(x, y).RGBA()
			nImg.Pix[offset] = uint8(r)
			nImg.Pix[offset+1] = uint8(g)
			nImg.Pix[offset+2] = uint8(b)
			nImg.Pix[offset+3] = uint8(a)

		}
	}

	return nImg, nil
}

// ImageToPixels transforms an image ot a pixel array
func ImageToPixels(img *image.NRGBA) []structures.Pixel {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels []structures.Pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := img.Stride*y + x*4

			pixels = append(pixels, structures.Pixel{R: img.Pix[offset], G: img.Pix[offset+1], B: img.Pix[offset+2], A: img.Pix[offset+3]})
		}
	}

	return pixels
}

// PixelsToImage transforms pixels to an image
func PixelsToImage(pixels []structures.Pixel, rect image.Rectangle) *image.NRGBA {
	img := image.NewNRGBA(rect)
	img.Stride = rect.Max.X * 4

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			offset := img.Stride*y + x*4
			img.Pix[offset] = pixels[y*rect.Max.X+x].R
			img.Pix[offset+1] = pixels[y*rect.Max.X+x].G
			img.Pix[offset+2] = pixels[y*rect.Max.X+x].B
			img.Pix[offset+3] = pixels[y*rect.Max.X+x].A
		}
	}

	return img
}

// SaveImage saves an image to a file
func SaveImage(img *image.NRGBA, name string) {

	fimg, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fimg.Close()

	png.Encode(fimg, img)

}
