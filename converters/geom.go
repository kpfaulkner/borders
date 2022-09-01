package converters

import (
	"fmt"
	"image"
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
func ConvertContourToPolygon(c *border.Contour, filterOutConflictingBoundaries bool, simplify bool, pointConverters ...PointConverter) (*geom.MultiPolygon, error) {

	//ls := []geom.LineString{}
	polygons := []geom.Polygon{}

	// err := convertContourToLineStrings(c, pointConverters, &ls, filterOutConflictingBoundaries)
	err := convertContourToPolygons(c, pointConverters, &polygons, filterOutConflictingBoundaries)
	if err != nil {
		return nil, err
	}

	mp, err := geom.NewMultiPolygon(polygons)
	if err != nil {
		fmt.Printf("Cannot make multipolygon: %s\n", err.Error())
		return nil, err
	}

	if simplify {

		// should calculate tolerance but seems way off. Keeping to 0.00015 for tests so far.
		//tolerance := generateSimplifyTolerance(22)
		tolerance := 0.00015
		// will calculate the threshold later. For now, 0.0002 is a reasonable value
		p2, err := mp.Simplify(tolerance, geom.DisableAllValidations)
		if err != nil {
			fmt.Printf("Cannot simplify polygon: %s\n", err.Error())
			return nil, err
		}
		return &p2, nil
	}
	return &mp, nil
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

func convertContourToPolygons(c *border.Contour, pointConverters []PointConverter, polygons *[]geom.Polygon, filterConflicts bool) error {

	// artificial bailout.

	// mark children with conflicts as unusable
	// any siblings that we conflict with, mark as unusable. This means that when multiple siblings (that conflict)
	// the first one will be used but others will not. Not ideal, but is a starting place. TODO(kpfaulkner) FIX THIS!
	//markConflictedSiblingsAsUnusable(c)

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
			err := convertContourToPolygons(child, pointConverters, polygons, filterConflicts)
			if err != nil {
				fmt.Printf("XXX err2 %s\n", err.Error())
				return err
			}
		}
	}

	return nil
}

// convertContourToLineStrings generate linestrings (later used to create polygons).
// filterConflicts is used to filter out boundaries that conflict with parent or siblings.
// Have NOT determined a good way to determine which sibling should be removed. TODO(kpfaulkner)
func convertContourToLineStrings(c *border.Contour, pointConverters []PointConverter, lineStrings *[]geom.LineString, filterConflicts bool) error {

	seq := pointsToSequence(c.Points, pointConverters)

	if c.BorderType == border.Hole {
		seq = seq.Reverse()
	}

	if seq.Length() > 2 {
		ls, err := geom.NewLineString(seq)
		if err != nil {
			fmt.Printf("seq len %d\n", seq.Length())
			return err
		}

		// if linestring only has 1 value, then ditch.
		if seq.Length() >= 1 {
			//fmt.Printf("LS is %s\n", ls.AsText())

			if c.BorderType == border.Hole {
				ls = ls.Reverse()
			}
			fmt.Printf("LS is %s\n", ls.AsText())

			*lineStrings = append(*lineStrings, ls)
		}
	}

	// mark children with conflicts as unusable
	// any siblings that we conflict with, mark as unusable. This means that when multiple siblings (that conflict)
	// the first one will be used but others will not. Not ideal, but is a starting place. TODO(kpfaulkner) FIX THIS!
	markConflictedSiblingsAsUnusable(c)

	for _, child := range c.Children {

		// only process child if no conflict with parent.
		if !child.ParentCollision && child.Usable {
			err := convertContourToLineStrings(child, pointConverters, lineStrings, filterConflicts)
			if err != nil {
				fmt.Printf("XXX err2 %s\n", err.Error())
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
	tolerance := simplifyToleranceInMetres / metresPerTile
	return tolerance
}

// tileSizeInMetres
func tileSizeInMetres(scale int) float64 {
	return 2 * math.Pi * EarthRadius / float64(uint64(1)<<uint64(scale))
}
