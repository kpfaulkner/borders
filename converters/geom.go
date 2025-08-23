package converters

import (
	"errors"
	"image"
	"math"

	"github.com/kpfaulkner/borders/border"
	"github.com/peterstace/simplefeatures/geom"
)

const (
	EarthRadius       = 6378137.0
	toleranceInMetres = 2

	radiansToDegreesRatio = math.Pi / 180.0
	degreesToRadiansRatio = 180.0 / math.Pi

	minLatitude  = -85.05112878
	maxLatitude  = 85.05112878
	minLongitude = -180
	maxLongitude = 180
)

type PointConverter func(x float64, y float64) (float64, float64)

// NewSlippyToLatLongConverter returns a function that converts slippy tile coordinates to lat/long.
func NewSlippyToLatLongConverter(slippyXOffset float64, slippyYOffset float64, scale int) func(X float64, Y float64) (float64, float64) {
	latLongN := math.Pow(2, float64(scale))
	f := func(x float64, y float64) (float64, float64) {
		long, lat := slippyCoordsToLongLat(slippyXOffset, slippyYOffset, x, y, latLongN)
		return long, lat
	}
	return f
}

// LatLongToSlippy converts lat/long to slippy tile coordinates.
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
// Simplifying while in "pixel space" simplifies the simplification degTolerance calculation.
// params:
//
//		simplify: Simplify the resulting polygons
//	    minPoints: Minimum number of vertices for a polygon to be considered valid. If less than this, then will be discarded. 0 means no minimum
//		degTolerance: Tolerance in pixels when simplifying. If set to 0, then will use defaults.
//		multiPolygonOnly: If the geometry is results in a GeometryCollection, then extract out the multipolygon part and return that.
//		pointConverters: Used to convert point co-ord systems. eg. slippy to lat/long.
func ConvertContourToPolygon(c *border.Contour, scale int, simplify bool, minPoints int, tolerance float64, multiPolygonOnly bool, pointConverters ...PointConverter) (*geom.Geometry, error) {
	polygons := []geom.Polygon{}

	err := convertContourToPolygons(c, minPoints, &polygons)
	if err != nil {
		return nil, err
	}

	mp := geom.NewMultiPolygon(polygons)

	if simplify {
		if tolerance == 0 {
			tolerance = generateSimplifyTolerance(scale)
		}
		gg := mp.AsGeometry()
		simplifiedGeom, err := gg.Simplify(tolerance, geom.NoValidate{})
		if err != nil {
			return nil, err
		}

		if multiPolygonOnly {
			if simplifiedGeom.Type() == geom.TypeMultiPolygon {
				mp, _ = simplifiedGeom.AsMultiPolygon()
				return returnConvertedGeometry(&mp, pointConverters...)
			}

			////////////////////////
			// Need to check if this is still possible.
			// need to check when we get geometrycollection vs multipolygon
			//if simplifiedGeom.Type() == geom.TypeGeometryCollection {
			//	gc, ok := simplifiedGeom.AsGeometryCollection()
			//	if ok {
			//		mp, err := filterMultiPolygonFromGeometryCollection(&gc)
			//		if err == nil {
			//			return returnConvertedGeometry(mp, pointConverters...)
			//		}
			//	}
			//}
			////////////////////////
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

// returnConvertedGeometry converts the multipolygon with PointConverters (if supplied)
// Can be used to help convert to lat/long or any other co-ordinate system.
func returnConvertedGeometry(mp *geom.MultiPolygon, pointConverters ...PointConverter) (*geom.Geometry, error) {
	finalMultiPoly, err := convertCoords(mp, pointConverters...)
	if err != nil {
		return nil, err
	}
	g := finalMultiPoly.AsGeometry()
	return &g, nil
}

// convertCoords converts the coordinates of a multipolygon using the supplied PointConverters.
func convertCoords(mp *geom.MultiPolygon, converters ...PointConverter) (*geom.MultiPolygon, error) {

	mp2 := mp.TransformXY(func(xy geom.XY) geom.XY {
		x := xy.X
		y := xy.Y
		// run through converters.
		for _, converter := range converters {
			newX, newY := converter(x, y)
			x = newX
			y = newY
		}
		return geom.XY{X: x, Y: y}
	})

	return &mp2, nil

}

// generateLineString generates a LineString from a slice of image.Points.
func generateLineString(points []image.Point) (*geom.LineString, error) {
	seq := pointsToSequence(points)

	if seq.Length() > 2 {
		ls := geom.NewLineString(seq)

		// if linestring only has 1 value, then ditch.
		if seq.Length() >= 1 {
			return &ls, nil
		}
	}

	return &geom.LineString{}, nil
}

// convertContourToPolygons converts the contour to a set of polygons but does NOT convert to different co-ord systems.
// If a polygon has fewer than minPoints then it will be discarded. 0 means no min points.
func convertContourToPolygons(c *border.Contour, minPoints int, polygons *[]geom.Polygon) error {

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
		if minPoints == 0 || len(lineStrings) > minPoints {
			poly = geom.NewPolygon(lineStrings)
			*polygons = append(*polygons, poly)
		}
	}

	for _, child := range c.Children {
		// only process child if no conflict with parent.
		if !child.ParentCollision && child.Usable {
			err := convertContourToPolygons(child, minPoints, polygons)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// pointsToSequence converts a slice of image.Points to a geom.Sequence.
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

// filterMultiPolygonFromGeometryCollection currently unused. Will be used in upcoming version.
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

// NewPixelToLatLongConverter returns a function that converts pixel coordinates to lat/long.
// Process is:
//
// 1) get X,Y coordinates for the topleft pixel
// 2) For each x,y coords passed (which will be position within image), convert to global space (add globalX/globalY)
// 3) Then run PixelXYToLatLong for each new globally positions pixel
func NewPixelToLatLongConverter(topLeftPixelLong float64, topLeftPixelLat float64, scale int) func(X float64, Y float64) (float64, float64) {

	// global pixel position of top left corner.
	gX, gY := LatLongToPixelXY(float64(topLeftPixelLat), float64(topLeftPixelLong), scale)
	globalX := float64(gX)
	globalY := float64(gY)
	f := func(x float64, y float64) (float64, float64) {
		newX := x + globalX
		newY := y + globalY
		lat, lon := PixelXYToLatLong(uint64(newX), uint64(newY), scale)
		return lon, lat
	}
	return f
}

func PixelXYToLatLong(pixelX uint64, pixelY uint64, scale int) (float64, float64) {

	// temp hack.. still trying to remember why
	//pixelX += 16
	//pixelY -= 1

	pixelTileSize := 256.0
	pixelGlobeSize := pixelTileSize * math.Pow(2, float64(scale))
	xPixelsToDegreesRatio := pixelGlobeSize / 360.0
	yPixelsToRadiansRatio := pixelGlobeSize / (2.0 * math.Pi)
	halfPixelGlobeSize := pixelGlobeSize / 2.0

	longitude := (float64(pixelX) - halfPixelGlobeSize) / xPixelsToDegreesRatio
	latitude := (2*math.Atan(math.Exp((float64(pixelY)-halfPixelGlobeSize)/(-yPixelsToRadiansRatio))) -
		math.Pi/2.0) * degreesToRadiansRatio

	return latitude, longitude
}

func LatLongToPixelXY(latitude float64, longitude float64, scale int) (uint64, uint64) {

	pixelTileSize := 256.0
	pixelGlobeSize := pixelTileSize * math.Pow(2, float64(scale))
	xPixelsToDegreesRatio := pixelGlobeSize / 360.0
	yPixelsToRadiansRatio := pixelGlobeSize / (2.0 * math.Pi)
	halfPixelGlobeSize := pixelGlobeSize / 2.0

	x := math.Round(halfPixelGlobeSize + (longitude * xPixelsToDegreesRatio))
	f := math.Min(math.Max(math.Sin(latitude*radiansToDegreesRatio), -0.9999), 0.9999)
	y := math.Round(halfPixelGlobeSize + 0.5*math.Log((1+f)/(1-f))*(-yPixelsToRadiansRatio))
	return uint64(x), uint64(y)

}
