package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"sync"
)

var exampleImg string = "./images/parisz.jpg"

func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
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

	// get image height and width with image/jpeg
	// change accordinly if file is png or gif

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

	// get the image
	img, _, err := image.Decode(imgfile)

	var wg sync.WaitGroup

	for y := 0; y < height; y++ {
		fmt.Println("Main: Starting worker", y)
		wg.Add(1)
		go printPixelsFor(&wg, img, y, width)
	}

	wg.Wait()
}

func printPixelsFor(wg *sync.WaitGroup, img image.Image, y int, width int) {
	defer wg.Done()
	for x := 0; x < width; x++ {
		r, g, b, a := img.At(x, y).RGBA()
		fmt.Printf("[X : %d Y : %v] R : %v, G : %v, B : %v, A : %v  \n", x, y, r, g, b, a)
	}

	fmt.Println("DONE!")
	os.Exit(1)
}
