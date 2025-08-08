package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kpfaulkner/borders/border"
)

func main() {
	PrintMemUsage("beginning")
	img, err := border.LoadImage("../../testimages/florida.png", 1, 1)
	if err != nil {
		panic("BOOM " + err.Error())
	}

	border.SaveImage("after-erode-dilate.png", img)

	PrintMemUsage("image loaded")

	start := time.Now()
	cont, err := border.FindContours(img)
	if err != nil {
		panic("BOOM " + err.Error())
	}
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	fmt.Printf("contour: %+v\n", cont.Children[0].Points)
	PrintMemUsage("found contours")
	border.SaveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0)
	PrintMemUsage("end")
}

// PrintMemUsage prints memory usage (allocation/GCs to stdout
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
