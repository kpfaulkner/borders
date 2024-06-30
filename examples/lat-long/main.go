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
	lng := 150.300446
	lat := -34.652429

	scale := 18
	img, err := border.LoadImage("../../testimages/highres-bw.png", 1, 1)

	start := time.Now()
	cont := border.FindContours(img)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	fmt.Printf("contour: %+v\n", cont.Children[0].Points)
	border.SaveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0)

	slippyX, slippyY := converters.LatLongToSlippy(lat, lng, scale)
	slippyConverter := converters.NewSlippyToLatLongConverter(slippyX, slippyY, scale)

	// tolerance of 0 means get ConvertContourToPolygon to calculate it
	poly, err := converters.ConvertContourToPolygon(cont, scale, true, 0, 0, true, slippyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to simple polygon : %s", err.Error())
	}

	j, _ := poly.MarshalJSON()
	os.WriteFile("final.geojson", j, 0644)

	fmt.Printf("convert to polygon took %d ms\n", time.Now().Sub(start).Milliseconds())
	b, _ := poly.MarshalJSON()
	fmt.Printf("%s\n", string(b))

}
