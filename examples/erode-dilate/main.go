package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kpfaulkner/borders/border"
	"github.com/kpfaulkner/borders/image"
)

func main() {
	PrintMemUsage("beginning")
	img, err := border.LoadImage("florida-big.png", false)

	// erode + dilate are used to remove noise from the image.
	img2, err := image.Erode(img, 1)
	if err != nil {
		panic("BOOM on erode")
	}

	img3, err := image.Dilate(img2, 1)
	if err != nil {
		panic("BOOM on dilate")
	}
	border.SaveImage("after-erode-dilate.png", img3)

	PrintMemUsage("image loaded")

	start := time.Now()
	cont := border.FindContours(img3)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	fmt.Printf("contour: %+v\n", cont.Children[0].Points)
	PrintMemUsage("found contours")
	border.SaveContourSliceImage("contour.png", cont, img3.Width, img3.Height, false, 0)
	PrintMemUsage("end")
}

func PrintMemUsage(header string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("=====  %s  =====\n", header)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
