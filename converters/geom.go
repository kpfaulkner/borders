package converters

import (
	"fmt"
	"image"
	"math"

	"github.com/kpfaulkner/borders/border"
	"github.com/peterstace/simplefeatures/geom"
)

type PointConverter func(x float64, y float64) (float64, float64)

func NewSlippyToLatLongConverter(slippyXOffset float64, slippyYOffset float64, scale int) func(X float64, Y float64) (float64, float64) {
	latLongN := math.Pow(2, float64(scale))
	f := func(x float64, y float64) (float64, float64) {
		long, lat := slippyCoordsToLongLat(slippyXOffset, slippyYOffset, x, y, latLongN)
		return long, lat
	}
	return f
}

// ConvertContourToPolygon converts the contours (set of x/y coords) to geometries commonly used in the GIS space
func ConvertContourToPolygon(c *border.Contour, pointConverters ...PointConverter) (*geom.Polygon, error) {

	ls := []geom.LineString{}

	err := convertContourToLineStrings(c, pointConverters, &ls)
	if err != nil {
		return nil, err
	}

	p, err := geom.NewPolygon(ls, geom.DisableAllValidations)
	if err != nil {
		return nil, err
	}

	// will calculate the threshold later. For now, 0.0002 is a reasonable value
	p2, err := p.Simplify(0.0002, geom.DisableAllValidations)
	if err != nil {
		return nil, err
	}

	return &p2, nil
}

// convertContourToLineStrings
func convertContourToLineStrings(c *border.Contour, pointConverters []PointConverter, lineStrings *[]geom.LineString) error {

	seq := pointsToSequence(c.Points, pointConverters)

	if seq.Length() > 2 {
		ls, err := geom.NewLineString(seq)
		if err != nil {
			fmt.Printf("seq len %d\n", seq.Length())
			return err
		}

		// if linestring only has 1 value, then ditch.
		if seq.Length() >= 1 {
			*lineStrings = append(*lineStrings, ls)
		}
	}
	for _, child := range c.Children {
		err := convertContourToLineStrings(child, pointConverters, lineStrings)
		if err != nil {
			fmt.Printf("XXX err2 %s\n", err.Error())
			return err
		}
	}

	return nil
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
