package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"sync"
	"time"
)

var exampleImg string = "./images/parisz.jpg"

func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	getAllPixelValues(false)
}

func getAllPixelValues(goroutines bool) {
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

	img, _, err := image.Decode(imgfile)

	if goroutines {
		var wg sync.WaitGroup
		start := time.Now()

		for y := 0; y < height; y++ {
			fmt.Println("Main: Starting worker", y)
			wg.Add(1)

			// go pixPrintPixelsForGoroutine(&wg, img, y, width)
			go printPixelsForGoroutine(&wg, img, y, width)
		}

		wg.Wait()
		duration := time.Since(start)
		fmt.Println("Operation (with goroutines) took:", duration)

	} else {

		start := time.Now()
		// pixPrintPixelsFor(img, height, width)
		printPixelsFor(img, height, width)

		duration := time.Since(start)
		fmt.Println("Operation took:", duration)
	}
}

func printPixelsForGoroutine(wg *sync.WaitGroup, img image.Image, y int, width int) {
	defer wg.Done()
	for x := 0; x < width; x++ {
		r, g, b, _ := img.At(x, y).RGBA()
		_ = uint8(r >> 8)
		_ = uint8(g >> 8)
		_ = uint8(b >> 8)

		//fmt.Printf("[X : %d Y : %v] R : %v, G : %v, B : %v, A : %v  \n", x, y, r, g, b, a)
	}

	fmt.Println("WORKER DONE.")
}

func pixPrintPixelsForGoroutine(wg *sync.WaitGroup, img image.Image, y int, width int) {
	defer wg.Done()
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	for x := 0; x < width; x++ {
		index := (y*width + x) * 4
		pix := rgba.Pix[index : index+4]
		_ = pix[0]
		_ = pix[1]
		_ = pix[2]

		//fmt.Printf("[X : %d Y : %v] R : %v, G : %v, B : %v \n", x, y, r, g, b)
	}

	fmt.Println("WORKER DONE.")
}

func printPixelsFor(img image.Image, height int, width int) {

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			_ = uint8(r >> 8)
			_ = uint8(g >> 8)
			_ = uint8(b >> 8)
		}
	}
}

func pixPrintPixelsFor(img image.Image, height int, width int) {
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
