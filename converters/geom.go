package converters

import (
	"fmt"
	"image"
	"math"

	"github.com/kpfaulkner/borders/border"
	"github.com/peterstace/simplefeatures/geom"
)

/*
type PointConverter interface {
	Convert(X float64, Y float64) (float64, float64)
} */

type PointConverter func(x float64, y float64) (float64, float64)

func NewSlippyToLatLongConverter(slippyXOffset float64, slippyYOffset float64, scale int) func(X float64, Y float64) (float64, float64) {
	latLongN := math.Pow(2, float64(scale))
	f := func(x float64, y float64) (float64, float64) {
		long, lat := slippyCoordsToLongLat(slippyXOffset, slippyYOffset, x, y, latLongN)
		return long, lat
	}
	return f
}

// ConvertContourToGeom converts the contours (set of x/y coords) to geometries commonly used in the GIS space
func ConvertContourToMultiPolygon(c *border.Contour, pointConverters ...PointConverter) (*geom.MultiPolygon, error) {
	root := geom.MultiPolygon{}

	_, _ = convertContourToPolygon(c, pointConverters)
	return &root, nil
}

// convertContourToPolygon converts the contours (set of x/y coords) to Polygon (and all children as well)
func convertContourToPolygon(c *border.Contour, pointConverters []PointConverter) (*geom.Polygon, error) {

	seq := pointsToSequence(c.Points, pointConverters)

	ls, err := geom.NewLineString(seq)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("XXXXXXXXX LS %s\n", ls.AsText())
	p, err := geom.NewPolygon([]geom.LineString{ls}, geom.DisableAllValidations)
	if err != nil {
		fmt.Printf("XXX err1 %s\n", err.Error())
		return nil, err
	}

	p2, err := p.Simplify(0.00002, geom.DisableAllValidations)
	if err != nil {
		fmt.Printf("BOOM simplify error %s\n", err.Error())
	}
	fmt.Printf("XXXXXX poly %s\n", p2.AsText())
	return nil, nil
}

func pointsToSequence(points []image.Point, converters []PointConverter) geom.Sequence {
	s := len(points)*2 + 2
	seq := make([]float64, s, s)
	index := 0
	for _, origP := range points {
		x, y := float64(origP.X), float64(origP.Y)

		// run through converters.
		for _, converter := range converters {
			newX, newY := converter(x, y)
			x, y = newX, newY
		}

		seq[index] = x
		seq[index+1] = y
		index += 2
	}

	seq[index] = seq[0]
	seq[index+1] = seq[1]
	return geom.NewSequence(seq, geom.DimXY)
}

// slippyCoordsToLongLat converts to lat/long... and requires the slippy offset of top left corner of area.
func slippyCoordsToLongLat(slippyXOffset float64, slippyYOffset float64, xTile float64, yTile float64, latLongN float64) (float64, float64) {
	//n := math.Pow(2, float64(scale))

	x := xTile + slippyXOffset
	y := yTile + slippyYOffset
	longDeg := (x/latLongN)*360.0 - 180.0
	latRad := math.Atan(math.Sinh(math.Pi - (y/latLongN)*2*math.Pi))
	latDeg := latRad * (180.0 / math.Pi)

	return longDeg, latDeg
}

// ConvertContourToMultiPolygonLatLong converts the contours (set of x/y coords) to geometries commonly used in the GIS space
// and converted to latitude/longitude
func ConvertContourToMultiPolygonLatLong(c *border.Contour, scale int, offsetX float64, offsetY float64) (*geom.MultiPolygon, error) {
	root := geom.MultiPolygon{}

	//_, _ = convertContourToPolygon(c)
	return &root, nil
}
