package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kpfaulkner/borders/border"
)

func displayContour(cont border.Contour) {

	for _, p := range cont.Points {
		fmt.Printf("(%d,%d)\n", p.X, p.Y)
	}

	for _, ch := range cont.Children {
		displayContour(*ch)
	}
}

func main() {
	PrintMemUsage("beginning")
	img, err := border.LoadImage("../../testimages/image3.png", 1, 1)
	PrintMemUsage("image loaded")
	if err != nil {
		panic("BOOM " + err.Error())
	}

	start := time.Now()
	cont := border.FindContours(img)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())
	displayContour(*cont)
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
