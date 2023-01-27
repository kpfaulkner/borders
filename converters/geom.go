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
	EarthRadius       = 6378137.0
	toleranceInMetres = 2
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

func LatLongToSlippy(latDegrees float64, longDegrees float64, scale int) (float64, float64) {
	n := math.Exp2(float64(scale))
	x := int(math.Floor((longDegrees + 180.0) / 360.0 * n))
	if float64(x) >= n {
		x = int(n - 1)
	}
	y := int(math.Floor((1.0 - math.Log(math.Tan(latDegrees*math.Pi/180.0)+1.0/math.Cos(latDegrees*math.Pi/180.0))/math.Pi) / 2.0 * n))
	return float64(x), float64(y)

}

// ConvertContourToPolygon converts the contours (set of x/y coords) to geometries commonly used in the GIS space
// Convert to polygons, then simplify (if required) while still in "pixel space"
// Only then apply conversions which may be to lat/long (or any other conversions).
// Simplifying while in "pixel space" simplifies the simplification tolerance calculation.
// params:
//
//	 	simplify. Simplify the resulting polygons
//		tolerance. Tolerance in pixels when simplifying. If set to 0, then will use defaults.
//		multiPolygonOnly. If the geometry is results in a GeometryCollection, then extract out the multipolygon part and return that.
//		pointConverters. Used to convert point co-ord systems. eg. slippy to lat/long.
func ConvertContourToPolygon(c *border.Contour, scale int, simplify bool, tolerance float64, multiPolygonOnly bool, pointConverters ...PointConverter) (*geom.Geometry, error) {
	polygons := []geom.Polygon{}

	err := convertContourToPolygons(c, &polygons)
	if err != nil {
		return nil, err
	}

	mp, err := geom.NewMultiPolygon(polygons)
	if err != nil {
		log.Errorf("Cannot make multipolygon: %s", err.Error())
		return nil, err
	}

	if simplify {
		if tolerance == 0 {
			tolerance = generateSimplifyTolerance(scale)
		}
		gg := mp.AsGeometry()
		simplifiedGeom, err := gg.Simplify(tolerance, geom.ConstructorOption(geom.DisableAllValidations))
		if err != nil {
			return nil, err
		}

		if multiPolygonOnly {
			if simplifiedGeom.Type() == geom.TypeMultiPolygon {
				mp, _ = simplifiedGeom.AsMultiPolygon()
				return returnConvertedGeometry(&mp, pointConverters...)
			}

			// need to check when we get geometrycollection vs multipolygon
			if simplifiedGeom.Type() == geom.TypeGeometryCollection {
				gc, ok := simplifiedGeom.AsGeometryCollection()
				if ok {
					mp, err := filterMultiPolygonFromGeometryCollection(&gc)
					if err == nil {
						return returnConvertedGeometry(mp, pointConverters...)
					}
				}
			}
			return nil, errors.New("unable to filter multipolygon from geometry collection")
		}
		mp, ok := simplifiedGeom.AsMultiPolygon()
		if ok {
			return returnConvertedGeometry(&mp, pointConverters...)
		} else {
			return nil, errors.New("unable to convert simplified geom to multipolygon")
		}
	}
	return returnConvertedGeometry(&mp, pointConverters...)
}

func returnConvertedGeometry(mp *geom.MultiPolygon, pointConverters ...PointConverter) (*geom.Geometry, error) {
	finalMultiPoly, err := convertCoords(mp, pointConverters...)
	if err != nil {
		return nil, err
	}
	g := finalMultiPoly.AsGeometry()
	return &g, nil
}

func convertCoords(mp *geom.MultiPolygon, converters ...PointConverter) (*geom.MultiPolygon, error) {

	mp2, err := mp.TransformXY(func(xy geom.XY) geom.XY {
		x := xy.X
		y := xy.Y
		// run through converters.
		for _, converter := range converters {
			newX, newY := converter(x, y)
			x = newX
			y = newY
		}
		return geom.XY{X: x, Y: y}
	}, geom.DisableAllValidations)

	if err != nil {
		log.Errorf("convertCoords err %s", err.Error())
	}
	return &mp2, err

}

func generateLineString(points []image.Point) (*geom.LineString, error) {
	seq := pointsToSequence(points)

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

// convertContourToPolygons converts the contour to a set of polygons but does NOT convert to different co-ord systems.
func convertContourToPolygons(c *border.Contour, polygons *[]geom.Polygon) error {

	// outer... so make a poly
	// will also cover hole if there.
	if c.BorderType == border.Outer {

		lineStrings := []geom.LineString{}
		outerLS, err := generateLineString(c.Points)
		if err != nil {
			return err
		}
		lineStrings = append(lineStrings, *outerLS)

		// now get children... (holes).
		for _, child := range c.Children {
			if !child.ParentCollision && child.Usable {
				ls, err := generateLineString(child.Points)
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
			err := convertContourToPolygons(child, polygons)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pointsToSequence(points []image.Point) geom.Sequence {
	s := len(points)*2 + 2
	seq := make([]float64, s, s)
	index := 0
	for _, origP := range points {
		x, y := float64(origP.X), float64(origP.Y)
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

// generateSimplifyTolerance will mainly be used when we want to convert to geographical co-ordinates
// By default we will determine how many metres per pixel (for input scale/zoom) and double it.
func generateSimplifyTolerance(scale int) float64 {
	mtrPerPixel := metresPerPixel(scale)
	tolerance := mtrPerPixel * toleranceInMetres
	return tolerance
}

// tileSizeInMetres is the size of a tile in metres.
func tileSizeInMetres(scale int) float64 {
	return 2 * math.Pi * EarthRadius / float64(uint64(1)<<uint64(scale))
}

// metresPerPixel is number of metres for a given input pixel. This is based on the scale/zoom.
func metresPerPixel(scale int) float64 {
	return tileSizeInMetres(scale) / 256.0
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

func NewPixelXYToLatLongConverter(latCentre float64, lonCentre float64, scale float64, imageWidth float64, imageHeight float64) func(X float64, Y float64) (float64, float64) {
	f := func(x float64, y float64) (float64, float64) {
		lat, lon := XYToLatLong(latCentre, lonCentre, int(scale), imageWidth, imageHeight, x, y)
		return lat, lon
	}
	return f
}

func XYToLatLong(latCentre float64, lonCentre float64, scale int, imageWidth float64, imageHeight float64, x float64, y float64) (float64, float64) {
	parallelMultiplier := math.Cos(latCentre * math.Pi / 180)
	degreesPerPixelX := 360 / math.Pow(2, float64(scale+8))
	degreesPerPixelY := 360 / math.Pow(2, float64(scale+8)) * parallelMultiplier
	pointLat := latCentre - degreesPerPixelY*(y-imageHeight/2)
	pointLng := lonCentre + degreesPerPixelX*(x-imageWidth/2)

	return pointLng, pointLat
}
