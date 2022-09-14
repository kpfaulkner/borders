package converters

import (
	"errors"
	"fmt"
	"image"
	"math"

	"github.com/kpfaulkner/borders/border"
	"github.com/peterstace/simplefeatures/geom"
	"github.com/peterstace/simplefeatures/geos"
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
func ConvertContourToPolygon(c *border.Contour, simplify bool, multiPolygonOnly bool, pointConverters ...PointConverter) (*geom.Geometry, error) {
	polygons := []geom.Polygon{}

	err := convertContourToPolygons(c, pointConverters, &polygons)
	if err != nil {
		return nil, err
	}

	mp, err := geom.NewMultiPolygon(polygons)
	if err != nil {
		fmt.Printf("Cannot make multipolygon: %s\n", err.Error())
		return nil, err
	}

	// need to find condition where MakeValid is actually required. Not optimal to do it all the time.
	gg, err := geos.MakeValid(mp.AsGeometry())
	if err != nil {
		return nil, err
	}

	if simplify {
		//tolerance := generateSimplifyTolerance(22)
		tolerance := 0.0002

		// will calculate the threshold later. For now, 0.0002 is a reasonable value
		simplifiedGeom, err := gg.Simplify(tolerance)
		if err != nil {
			return nil, err
		}

		if multiPolygonOnly {
			gc, ok := simplifiedGeom.AsGeometryCollection()
			if ok {
				mp, err := filterMultiPolygonFromGeometryCollection(&gc)
				if err == nil {
					g2 := mp.AsGeometry()
					return &g2, nil
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
			fmt.Printf("seq len %d\n", seq.Length())
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
			fmt.Printf("unable to make polygon, len %d : %s\n", len(lineStrings), err.Error())
			poly, err = geom.NewPolygon(lineStrings)
			if err != nil {
				fmt.Printf("unable to make polygon second time : %s\n", err.Error())
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

func markConflictedSiblingsAsUnusable(node *border.Contour) {

	conflictedIds := make(map[int]bool)
	for _, c := range node.Children {

		// if already marked as unusable, then skip
		if !c.Usable {
			continue
		}

		// collision with parent... skip
		if c.ParentCollision {
			//fmt.Printf("1 marking %d as unusable\n", c.Id)
			//c.Usable = false
			continue
		}

		// already marked in conflictedIds... skip
		if _, has := conflictedIds[c.Id]; has {
			//fmt.Printf("2 marking %d as unusable\n", c.Id)
			//c.Usable = false
		}

		// take all the ones that this child is conflicting with, and mark those to skip.
		for conflictId, _ := range c.ConflictingContours {
			conflictedIds[conflictId] = true
		}
	}
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
	metresPerTile := tileSizeInMetres(scale)
	fmt.Printf("metres per tile %f\n", metresPerTile)
	tolerance := simplifyToleranceInMetres / metresPerTile
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
