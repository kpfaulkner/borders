package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kpfaulkner/borders/border"
	"github.com/kpfaulkner/borders/converters"
)

func main() {
	fmt.Printf("So it begins...\n")

	//img, err := border.LoadImage("image1.png")
	//img, err := border.LoadImage("image3.png")
	//img, err := border.LoadImage("tiny.png")
	//img, err := border.LoadImage("big-test-image.png")
	img, err := border.LoadImage("big-image2.png")

	if err != nil {
		panic("BOOM " + err.Error())
	}

	start := time.Now()
	cont := border.FindContours(img)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	//saveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0, false)
	border.SaveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0)
	//border.SaveContourSliceImage("c:/temp/contour/contour", cont, img.Width, img.Height, true, 0)

	slippyConverter := converters.NewSlippyToLatLongConverter(1139408, 1772861, 22)

	poly, err := converters.ConvertContourToPolygon(cont, slippyConverter)
	if err != nil {
		log.Fatalf("Unable to convert to polygon : %s", err.Error())
	}

	b, _ := poly.MarshalJSON()
	fmt.Printf(string(b))

}
