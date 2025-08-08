package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kpfaulkner/borders/border"
	"github.com/kpfaulkner/borders/converters"
)

func main() {
	scale := 21
	img, err := border.LoadImage("../../testimages/testmap2-1891519-1285047-21.png", 1, 1)

	start := time.Now()
	cont, err := border.FindContours(img)
	if err != nil {
		log.Fatalf("Unable to find contours: %s", err.Error())
		return
	}
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	// save the contour as an image.
	border.SaveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0)

	// create a slippy map converter so we can generate a multipolygon that has the latitude/longitude
	// geojson.
	// This means that we need to know the slippy X/Y coords of the top left corner of the map
	slippyX := 1891519.0
	slippyY := 1285047.0
	slippyConverter := converters.NewSlippyToLatLongConverter(slippyX, slippyY, scale)

	// tolerance of 0 means get ConvertContourToPolygon to calculate it
	poly, err := converters.ConvertContourToPolygon(cont, scale, true, 0, 0, true, slippyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to simple polygon : %s", err.Error())
	}

	j, _ := poly.MarshalJSON()
	os.WriteFile("final.geojson", j, 0644)

	fmt.Printf("convert to polygon took %d ms\n", time.Now().Sub(start).Milliseconds())
}
