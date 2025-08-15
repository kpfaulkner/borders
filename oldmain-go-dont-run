package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/kpfaulkner/borders/border"
	"github.com/kpfaulkner/borders/converters"
)

func main() {
	fmt.Printf("So it begins...\n")

	// example lat/lon of top left of testmap2-1891519-1285047-21.png
	lng := 144.700927662535
	lat := -37.5696166965532
	scale := 21

	slippyX, slippyY := converters.LatLongToSlippy(lat, lng, scale)
	img, err := border.LoadImage(`c:\temp\testcase2.png`, 1, 1)

	border.SaveImage("bordertest.png", img)

	PrintMemUsage("image loaded")
	if err != nil {
		panic("BOOM " + err.Error())
	}

	start := time.Now()
	cont, err := border.FindContours(img)
	if err != nil {
		fmt.Printf("Unable to find contours : %s", err.Error())
		return
	}
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	fmt.Printf("contour: %+v\n", cont.Children[0].Points)
	PrintMemUsage("found contours")
	border.SaveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0)

	// If the input image are base off a slippy mask (ie each pixel represents a tile) then we require a slippy converter.
	slippyConverter := converters.NewSlippyToLatLongConverter(slippyX, slippyY, scale)

	// tolerance of 0 means get ConvertContourToPolygon to calculate it
	poly, err := converters.ConvertContourToPolygon(cont, scale, true, 0, 0, true, slippyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to simple polygon : %s", err.Error())
	}

	j, _ := poly.MarshalJSON()
	os.WriteFile("simple-final.geojson", j, 0644)

	fmt.Printf("convert to simple polygon took %d ms\n", time.Now().Sub(start).Milliseconds())

	// tolerance of 0 means get ConvertContourToPolygon to calculate it
	fullPoly, err := converters.ConvertContourToPolygon(cont, scale, false, 0, 0, true, slippyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to polygon : %s", err.Error())
	}

	j, _ = fullPoly.MarshalJSON()
	os.WriteFile("full-final.geojson", j, 0644)

	fmt.Printf("convert to polygon took %d ms\n", time.Now().Sub(start).Milliseconds())
	PrintMemUsage("convert to poly")

	b, _ := poly.MarshalJSON()
	fmt.Printf("%s\n", string(b))

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
