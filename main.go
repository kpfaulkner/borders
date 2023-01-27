package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/kpfaulkner/borders/border"
	"github.com/kpfaulkner/borders/converters"
	"github.com/kpfaulkner/borders/image"
)

func main() {
	fmt.Printf("So it begins...\n")

	lng := 150.300446
	lat := -34.652429
	scale := 18

	slippyX, slippyY := converters.LatLongToSlippy(lat, lng, scale)
	img, err := border.LoadImage("image based off slippy mask", false)

	img2, err := image.Erode(img, 1)
	if err != nil {
		panic("BOOM on erode")
	}

	img3, err := image.Dilate(img2, 1)
	if err != nil {
		panic("BOOM on dilate")
	}
	border.SaveImage("bordertest.png", img3)

	PrintMemUsage("image loaded")
	if err != nil {
		panic("BOOM " + err.Error())
	}

	start := time.Now()
	cont := border.FindContours(img3)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	fmt.Printf("contour: %+v\n", cont.Children[0].Points)
	PrintMemUsage("found contours")
	border.SaveContourSliceImage("contour.png", cont, img3.Width, img3.Height, false, 0)

	// If the input image are base off a slippy mask (ie each pixel represents a tile) then we require a slippy converter.
	slippyConverter := converters.NewSlippyToLatLongConverter(slippyX, slippyY, scale)

	// tolerance of 0 means get ConvertContourToPolygon to calculate it
	poly, err := converters.ConvertContourToPolygon(cont, scale, true, 0, true, slippyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to polygon : %s", err.Error())
	}

	j, _ := poly.MarshalJSON()
	os.WriteFile("final.geojson", j, 0644)

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
