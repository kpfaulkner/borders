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
	img, err := border.LoadImage("../../testimages/airport.png", 1, 1)

	start := time.Now()
	cont, err := border.FindContours(img)
	if err != nil {
		log.Fatalf("Unable to find contours: %s", err.Error())
		return
	}
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	// save the contour as an image.
	border.SaveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0)

	// lat/lon of top left of image.
	lng := 144.843699326
	lat := -37.667085056

	// Converter from pixels to lat/lon
	conv := converters.NewPixelToLatLongConverter(lng, lat, scale)

	// tolerance of 0 means get ConvertContourToPolygon to calculate it
	poly, err := converters.ConvertContourToPolygon(cont, scale, false, 0, 0, true, conv)
	if err != nil {
		log.Fatalf("Unable to convert to simple polygon : %s", err.Error())
	}

	j, _ := poly.MarshalJSON()
	os.WriteFile("final.geojson", j, 0644)

	fmt.Printf("convert to polygon took %d ms\n", time.Now().Sub(start).Milliseconds())
}
