package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"time"
)

var exampleImg string = "./images/parisz.jpg"

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	getAllPixelValues()
}

func getAllPixelValues() {
	imgfile, err := os.Open(exampleImg)
	if err != nil {
		fmt.Printf("location '%v' not found!", exampleImg)
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
	getImageCopyOf(img, false, height, width)

	duration := time.Since(start)
	fmt.Println("Operation took:", duration)
}

func getImageCopyOf(img image.Image, imgMethod bool, height int, width int) {
	// use img.At() method of extracting RGB from pixels in given image
	if imgMethod {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				_ = uint8(r >> 8)
				_ = uint8(g >> 8)
				_ = uint8(b >> 8)
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
				_ = pix[0]
				_ = pix[1]
				_ = pix[2]
			}
		}
	}
}
