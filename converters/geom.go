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

	degreesToRadiansRatio = math.Pi / 180.0
	radiansToDegreesRatio = 180.0 / math.Pi
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
				return returnConvertedGeometry(&mp, pointConverters...), nil
			}
			return nil, errors.New("unable to filter multipolygon from geometry collection")
		}

		mp, ok := simplifiedGeom.AsMultiPolygon()
		if ok {
			return returnConvertedGeometry(&mp, pointConverters...), nil
		} else {
			return nil, errors.New("unable to convert simplified geom to multipolygon")
		}
	}
	return returnConvertedGeometry(&mp, pointConverters...), nil
}

// returnConvertedGeometry converts the multipolygon with PointConverters (if supplied)
// Can be used to help convert to lat/long or any other co-ordinate system.
func returnConvertedGeometry(mp *geom.MultiPolygon, pointConverters ...PointConverter) *geom.Geometry {
	finalMultiPoly := convertCoords(mp, pointConverters...)
	g := finalMultiPoly.AsGeometry()
	return &g
}

// convertCoords converts the coordinates of a multipolygon using the supplied PointConverters.
func convertCoords(mp *geom.MultiPolygon, converters ...PointConverter) *geom.MultiPolygon {
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
	return &mp2
}

// generateLineString generates a LineString from a slice of image.Points.
func generateLineString(points []image.Point) (*geom.LineString, error) {
	seq := pointsToSequence(points)

	if seq.Length() > 2 {
		ls := geom.NewLineString(seq)
		return &ls, nil
	}

	return &geom.LineString{}, nil
}

// convertContourToPolygons converts the contour to a set of polygons but does NOT convert to different co-ord systems.
// A polygon is discarded when its outer ring has fewer than minPoints vertices. minPoints == 0 disables the filter.
// Contours with an empty outer ring are always skipped, and empty hole rings are dropped so they don't poison the polygon.
func convertContourToPolygons(c *border.Contour, minPoints int, polygons *[]geom.Polygon) error {
	if c.BorderType == border.Outer && len(c.Points) > 0 && (minPoints == 0 || len(c.Points) >= minPoints) {

		lineStrings := []geom.LineString{}
		outerLS, err := generateLineString(c.Points)
		if err != nil {
			return err
		}
		lineStrings = append(lineStrings, *outerLS)

		// holes — skip any with no points, which would otherwise be appended as an empty LineString.
		for _, child := range c.Children {
			if !child.ParentCollision && child.Usable && len(child.Points) > 0 {
				ls, err := generateLineString(child.Points)
				if err != nil {
					return err
				}
				lineStrings = append(lineStrings, *ls)
			}
		}

		poly := geom.NewPolygon(lineStrings)
		*polygons = append(*polygons, poly)
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
// topLeftPixelLat/topLeftPixelLong are the geographic coordinates of the image's top-left pixel.
// Important to note that the converter function returned when executed will return
// (longitude,latitude) in that order.
// Process is:
//
// 1) get X,Y coordinates for the topleft pixel
// 2) For each x,y coords passed (which will be position within image), convert to global space (add globalX/globalY)
// 3) Convert each globally positioned pixel to lat/long using the precomputed scale-dependent constants.
func NewPixelToLatLongConverter(topLeftPixelLat float64, topLeftPixelLong float64, scale int) func(X float64, Y float64) (float64, float64) {

	// Scale-dependent constants — computed once here rather than per call.
	// PixelXYToLatLong recomputes these every invocation; the closure below is
	// called once per polygon vertex (often tens of thousands of times per
	// conversion), so hoisting the math.Exp2 and three divisions matters.
	pixelGlobeSize := 256.0 * math.Exp2(float64(scale))
	halfPixelGlobeSize := pixelGlobeSize / 2.0
	xPixelsToDegreesRatio := pixelGlobeSize / 360.0
	yPixelsToRadiansRatio := pixelGlobeSize / (2.0 * math.Pi)

	// global pixel position of top left corner.
	gX, gY := LatLongToPixelXY(topLeftPixelLat, topLeftPixelLong, scale)
	globalX := float64(gX)
	globalY := float64(gY)

	return func(x float64, y float64) (float64, float64) {
		pixelX := x + globalX
		pixelY := y + globalY
		longitude := (pixelX - halfPixelGlobeSize) / xPixelsToDegreesRatio
		latitude := (2*math.Atan(math.Exp((pixelY-halfPixelGlobeSize)/(-yPixelsToRadiansRatio))) - math.Pi/2.0) * radiansToDegreesRatio
		return longitude, latitude
	}
}

func PixelXYToLatLong(pixelX int64, pixelY int64, scale int) (float64, float64) {

	pixelTileSize := 256.0
	pixelGlobeSize := pixelTileSize * math.Pow(2, float64(scale))
	xPixelsToDegreesRatio := pixelGlobeSize / 360.0
	yPixelsToRadiansRatio := pixelGlobeSize / (2.0 * math.Pi)
	halfPixelGlobeSize := pixelGlobeSize / 2.0

	longitude := (float64(pixelX) - halfPixelGlobeSize) / xPixelsToDegreesRatio
	latitude := (2*math.Atan(math.Exp((float64(pixelY)-halfPixelGlobeSize)/(-yPixelsToRadiansRatio))) -
		math.Pi/2.0) * radiansToDegreesRatio

	return latitude, longitude
}

func LatLongToPixelXY(latitude float64, longitude float64, scale int) (int64, int64) {

	pixelTileSize := 256.0
	pixelGlobeSize := pixelTileSize * math.Pow(2, float64(scale))
	xPixelsToDegreesRatio := pixelGlobeSize / 360.0
	yPixelsToRadiansRatio := pixelGlobeSize / (2.0 * math.Pi)
	halfPixelGlobeSize := pixelGlobeSize / 2.0

	x := math.Round(halfPixelGlobeSize + (longitude * xPixelsToDegreesRatio))
	f := math.Min(math.Max(math.Sin(latitude*degreesToRadiansRatio), -0.9999), 0.9999)
	y := math.Round(halfPixelGlobeSize + 0.5*math.Log((1+f)/(1-f))*(-yPixelsToRadiansRatio))
	return int64(x), int64(y)

}
