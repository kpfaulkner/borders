package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/kpfaulkner/borders/border"
	"github.com/kpfaulkner/borders/converters"
	"github.com/kpfaulkner/borders/image"
)

func main() {
	fmt.Printf("So it begins...\n")

	PrintMemUsage("beginning")
	//img, err := border.LoadImage("testimages/image1.png")
	//img, err := border.LoadImage("testimages/image3.png")
	//img, err := border.LoadImage("testimages/image-simple.png")
	//img, err := border.LoadImage("testimages/image-simple3.png")
	//img, err := border.LoadImage("testimages/image-simple4.png")
	//img, err := border.LoadImage("testimages/image-simple5.png")
	//img, err := border.LoadImage("testimages/image-simple-thin-line.png")
	//img, err := border.LoadImage("tiny.png")
	//img, err := border.LoadImage("big-test-image.png")
	//img, err := border.LoadImage("testimages/sidespike.png", true)
	img, err := border.LoadImage("florida-big.png", false)
	//img, err := border.LoadImage("test-full.png", false)

	img2, err := image.Erode(img, 1)
	if err != nil {
		panic("BOOM on erode")
	}

	img3, err := image.Dilate(img2, 1)
	if err != nil {
		panic("BOOM on dilate")
	}

	//img3 := img2

	border.SaveImage("test.png", img3)

	PrintMemUsage("image loaded")
	if err != nil {
		panic("BOOM " + err.Error())
	}

	start := time.Now()
	cont := border.FindContours(img3)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	fmt.Printf("contour: %+v\n", cont.Children[0].Points)
	PrintMemUsage("found contours")
	//saveContourSliceImage("contour.png", cont, img.Width,
	//img.Height, false, 0, false)
	border.SaveContourSliceImage("contour.png", cont, img3.Width, img3.Height, false, 0)
	//border.SaveContourSliceImage("c:/temp/contour/contour", cont, img.Width, img.Height, true, 0)

	slippyConverter := converters.NewSlippyToLatLongConverter(1139408, 1772861, 22)

	poly, err := converters.ConvertContourToPolygon(cont, true, true, slippyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to polygon : %s", err.Error())
	}

	fmt.Printf("convert to polygon took %d ms\n", time.Now().Sub(start).Milliseconds())
	PrintMemUsage("convert to poly")

	b, _ := poly.MarshalJSON()
	fmt.Printf("%s\n", string(b))

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
