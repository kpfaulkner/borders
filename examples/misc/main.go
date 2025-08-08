package main

import (
	"fmt"

	"github.com/kpfaulkner/borders/converters"
)

func main() {
	lng := 151.09577747524995
	lat := -33.92072681353907
	scale := 20

	x, y := converters.LatLongToSlippy(lat, lng, scale)

	conv := converters.NewSlippyToLatLongConverter(x, y, scale)
	newLon, newLat := conv(x, y)

	fmt.Printf("orig long %f lat %f\n", lng, lat)
	fmt.Printf("slippy %f %f\n", x, y)
	fmt.Printf("new long %f lat %f\n", newLon, newLat)
}
