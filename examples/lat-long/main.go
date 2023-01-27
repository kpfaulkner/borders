package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kpfaulkner/borders/border"
	"github.com/kpfaulkner/borders/converters"
	"github.com/kpfaulkner/borders/image"
)

func main() {
	lng := 150.300446
	lat := -34.652429

	scale := 18
	img, err := border.LoadImage("../../testimages/highres-bw.png", false)

	img2, err := image.Erode(img, 1)
	if err != nil {
		panic("BOOM on erode")
	}

	img3, err := image.Dilate(img2, 1)
	if err != nil {
		panic("BOOM on dilate")
	}

	if err != nil {
		panic("BOOM " + err.Error())
	}

	start := time.Now()
	cont := border.FindContours(img3)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	fmt.Printf("contour: %+v\n", cont.Children[0].Points)
	border.SaveContourSliceImage("contour.png", cont, img3.Width, img3.Height, false, 0)

	xyConverter := converters.NewPixelXYToLatLongConverter(lat, lng, float64(scale), float64(img3.Width), float64(img3.Height))

	poly, err := converters.ConvertContourToPolygon(cont, scale, true, 0, true, xyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to polygon : %s", err.Error())
	}

	j, _ := poly.MarshalJSON()
	os.WriteFile("final.geojson", j, 0644)

	fmt.Printf("convert to polygon took %d ms\n", time.Now().Sub(start).Milliseconds())
	b, _ := poly.MarshalJSON()
	fmt.Printf("%s\n", string(b))

}
