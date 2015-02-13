package main

import (
	//"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	//"io"
	//"math"
	"os"
	"path/filepath"

	_ "code.google.com/p/vp8-go/webp"
	_ "image/jpeg"
)

var smallestWidth = 10000
var smallestHeight = 10000
var imageCount = 0

var jobs chan string
var results chan error
var up chan bool

func worker(id int, jobs <-chan string, results chan<- error, up chan<- bool) {
	fmt.Println("Started worker thread:", id)
	up <- true
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		imgFile, err := os.Open(j)
		if err != nil {
			results <- err
			return
		}
		img, _, err := image.Decode(imgFile)
		if err != nil {
			results <- err
			return
		}
		img = resize.Resize(uint(smallestWidth), uint(smallestHeight), img, resize.Lanczos3)
		w, err := os.Create(j + "_resized.png")
		if err != nil {
			results <- err
			return
		}
		err = png.Encode(w, img)
		w.Close()
		if err != nil {
			results <- err
			return
		}
		if err != nil {
			results <- err
			return
		}
	}
}

func main() {
	filepath.Walk(os.Args[1], WalkFunc)
	fmt.Printf("smallest image %d %d\n", smallestWidth, smallestHeight)

	jobs = make(chan string, imageCount)
	results = make(chan error, imageCount)
	up = make(chan bool, 10)

	// This starts up 10 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 10; w++ {
		go worker(w, jobs, results, up)
	}
	for w := 1; w <= 10; w++ {
		<-up
	}

	filepath.Walk(os.Args[1], WalkFunc2)
	close(jobs)

	// Finally we collect all the results of the work.
	for i := 0; i < imageCount; i++ {
		err := <-results
		fmt.Println("Error resizing:", err)
	}
}

func WalkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		imageCount++
		imgFile, _ := os.Open(path)
		img, _, _ := image.Decode(imgFile)

		if img.Bounds().Dx() < smallestWidth {
			smallestWidth = img.Bounds().Dx()
		}

		if img.Bounds().Dy() < smallestHeight {
			smallestHeight = img.Bounds().Dy()
		}
	}
	return nil
}

func WalkFunc2(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		fmt.Println("Sending", path, "to worker pool")
		jobs <- path
	}
	return nil
}
