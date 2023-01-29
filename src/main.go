package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"time"
)

var inputImg string = "./images/parisz.jpg"
var outputImg string = "images/copyimg.png"

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	getAllPixelValues()
}

func getAllPixelValues() {
	imgfile, err := os.Open(inputImg)
	if err != nil {
		fmt.Printf("location '%v' not found!", inputImg)
		os.Exit(1)
	}
	defer imgfile.Close()

	imgCfg, _, err := image.DecodeConfig(imgfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)

	// we need to reset the io.Reader again for image.Decode() function below to work
	// otherwise we will  - panic: runtime error: invalid memory address or nil pointer dereference
	// there is no build in rewind for io.Reader, use Seek(0,0)
	imgfile.Seek(0, 0)
	img, _, err := image.Decode(imgfile)

	start := time.Now()
	createImageCopyOf(img, false, height, width)

	duration := time.Since(start)
	fmt.Println("Operation took:", duration)
}

func createImageCopyOf(img image.Image, imgMethod bool, height int, width int) {
	copyImg := getImageBuilderTemplate(height, width)

	// use img.At() method of extracting RGB from pixels in given image
	if imgMethod {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				r8 := uint8(r >> 8)
				g8 := uint8(g >> 8)
				b8 := uint8(b >> 8)

				currentPixelColor := color.RGBA{r8, g8, b8, 0xff}
				copyImg.Set(x, y, currentPixelColor)
			}
		}

		// use rgba.Pix() method of extracting RGB from pixels in given image
	} else {
		bounds := img.Bounds()
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				index := (y*width + x) * 4
				pix := rgba.Pix[index : index+4]
				r8 := pix[0]
				g8 := pix[1]
				b8 := pix[2]

				currentPixelColor := color.RGBA{r8, g8, b8, 0xff}
				copyImg.Set(x, y, currentPixelColor)
			}
		}
	}

	createPNGEncodedImageFrom(copyImg)
}

// Return clean image template of specified size.
// Used to manipulate pixels on.
func getImageBuilderTemplate(height int, width int) *image.RGBA {
	upLeft := image.Point{0, 0}
	downRight := image.Point{width, height}
	return image.NewRGBA(image.Rectangle{upLeft, downRight})
}

// Create and save .png image from specified *image.RGBA.
func createPNGEncodedImageFrom(copyImg *image.RGBA) {
	f, err := os.Create("images/copyImg.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Successfully copied from '%v' to '%v'.\n", inputImg, outputImg)
	png.Encode(f, copyImg)
}
