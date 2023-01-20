package converters

import (
	"errors"
	"image"

	log "github.com/sirupsen/logrus"

	"math"

	"github.com/kpfaulkner/borders/border"
	"github.com/peterstace/simplefeatures/geom"
)

const (
	EarthRadius               = 6378137.0
	simplifyToleranceInMetres = 50.0
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
// Will most likely (after simplification) generate a GeometryCollection. Will need to strip out multipolygons from that.
// params:
//
//	simplify. Simplify the resulting polygons
//	multiPolygonOnly. If the geometry is results in a GeometryCollection, then extract out the multipolygon part and return that.
//	pointConverters. Used to convert point co-ord systems. eg. slippy to lat/long.
func ConvertContourToPolygon(c *border.Contour, scale int, simplify bool, multiPolygonOnly bool, pointConverters ...PointConverter) (*geom.Geometry, error) {
	polygons := []geom.Polygon{}

	err := convertContourToPolygons(c, pointConverters, &polygons)
	if err != nil {
		return nil, err
	}

	mp, err := geom.NewMultiPolygon(polygons)
	if err != nil {
		log.Errorf("Cannot make multipolygon: %s", err.Error())
		return nil, err
	}

	gg := mp.AsGeometry()
	if simplify {
		tolerance := generateSimplifyTolerance(scale)

		// will calculate the threshold later. For now, 0.0002 is a reasonable value
		simplifiedGeom, err := gg.Simplify(tolerance, geom.ConstructorOption(geom.DisableAllValidations))
		if err != nil {
			return nil, err
		}

		if multiPolygonOnly {
			if simplifiedGeom.Type() == geom.TypeMultiPolygon {
				return &simplifiedGeom, nil
			}

			// need to check when we get geometrycollection vs multipolygon
			if simplifiedGeom.Type() == geom.TypeGeometryCollection {
				gc, ok := simplifiedGeom.AsGeometryCollection()
				if ok {
					mp, err := filterMultiPolygonFromGeometryCollection(&gc)
					if err == nil {
						g2 := mp.AsGeometry()
						return &g2, nil
					}
				}
			}
			return nil, errors.New("unable to filter multipolygon from geometry collection")
		}
		return &simplifiedGeom, nil
	}
	return &gg, nil
}

func generateLineString(points []image.Point, pointConverters []PointConverter) (*geom.LineString, error) {
	seq := pointsToSequence(points, pointConverters)

	if seq.Length() > 2 {
		ls, err := geom.NewLineString(seq)
		if err != nil {
			return nil, err
		}

		// if linestring only has 1 value, then ditch.
		if seq.Length() >= 1 {
			return &ls, nil
		}
	}

	return &geom.LineString{}, nil
}

func convertContourToPolygons(c *border.Contour, pointConverters []PointConverter, polygons *[]geom.Polygon) error {

	// outer... so make a poly
	// will also cover hole if there.
	if c.BorderType == border.Outer {

		lineStrings := []geom.LineString{}
		outerLS, err := generateLineString(c.Points, pointConverters)
		if err != nil {
			return err
		}

		//*outerLS = outerLS.Simplify(0.0002)
		lineStrings = append(lineStrings, *outerLS)

		// now get children... (holes).
		for _, child := range c.Children {
			if !child.ParentCollision && child.Usable {
				ls, err := generateLineString(child.Points, pointConverters)
				if err != nil {
					return err
				}
				lineStrings = append(lineStrings, *ls)
			}
		}

		var poly geom.Polygon
		poly, err = geom.NewPolygon(lineStrings, geom.DisableAllValidations)
		if err != nil {
			log.Debugf("unable to make polygon, len %d : %s", len(lineStrings), err.Error())
			poly, err = geom.NewPolygon(lineStrings)
			if err != nil {
				log.Errorf("unable to make polygon second time : %s\n", err.Error())
				return err
			}
		}
		*polygons = append(*polygons, poly)
	}

	for _, child := range c.Children {
		// only process child if no conflict with parent.
		if !child.ParentCollision && child.Usable {
			err := convertContourToPolygons(child, pointConverters, polygons)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pointsToSequence(points []image.Point, converters []PointConverter) geom.Sequence {
	s := len(points)*2 + 2
	seq := make([]float64, s, s)
	index := 0
	for _, origP := range points {
		x, y := float64(origP.X)-0.5, float64(origP.Y)-0.5 // based off testing. Need to check WHY!??!

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

func generateSimplifyTolerance(scale int) float64 {

	// need to figure this out! TODO(kpfaulkner)
	/*
		metresPerTile := tileSizeInMetres(scale)
		tolerance := simplifyToleranceInMetres / metresPerTile
	*/

	// hardcode 0.0002 for now... working well for all test cases
	tolerance := 0.0002

	return tolerance
}

// tileSizeInMetres
func tileSizeInMetres(scale int) float64 {
	return 2 * math.Pi * EarthRadius / float64(uint64(1)<<uint64(scale))
}

func filterMultiPolygonFromGeometryCollection(col *geom.GeometryCollection) (*geom.MultiPolygon, error) {
	var mp geom.MultiPolygon
	var ok bool
	for i := 0; i < col.NumGeometries(); i++ {
		g := col.GeometryN(i)
		mp, ok = g.AsMultiPolygon()
		if ok {
			return &mp, nil
		}
	}

	return nil, errors.New("no multipolygon found in geometry collection")
}
