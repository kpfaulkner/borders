package border

import (
	"fmt"
	"image"
	"strings"
)

const (
	Outer = 0
	Hole  = 1
)

// Contour represents a single contour/border extracted from an image.
// It also tracks its parents and children.
type Contour struct {

	// Points making up the contour
	Points []image.Point

	Id int

	// Outer or Hole.
	BorderType int

	// Id of parent
	ParentId int

	// ParentCollision indicates if colliding with parent. Just an optimisation for quick removal later on.
	ParentCollision bool

	// Parent links to contours parent
	Parent *Contour

	// Children links to contours children
	Children []*Contour

	// ConflictingContours is a map of contours that we KNOW we conflict with. This may be the parent or other
	// siblings
	ConflictingContours map[int]bool // hate to use maps here... but want uniqueness
}

// NewContour create new contour
func NewContour(id int) *Contour {
	c := Contour{}
	c.Id = id
	c.BorderType = Hole
	c.ConflictingContours = make(map[int]bool)
	return &c
}

// AddPoint adds a point (image.Point) to the contour
func (c *Contour) AddPoint(p image.Point) error {
	c.Points = append(c.Points, p)
	return nil
}

// ContourStats generates writes to stdout stats about the contour and all children.
// Primarily used for debugging
func ContourStats(c *Contour, offset int) {
	if len(c.Points) > 0 {
		pad := strings.Repeat(" ", offset)
		fmt.Printf("%s%d : len %d :  no kids %d : no col %d : col with parent %+v\n", pad, c.Id, len(c.Points), len(c.Children), len(c.ConflictingContours), c.ParentCollision)
	}

	for _, ch := range c.Children {
		ContourStats(ch, offset+2)
	}
}

// ContourStatsWithCollisions generates writes to stdout stats about the contour and all children that have collisions
// Primarily used for debugging.
func ContourStatsWithCollisions(c *Contour, offset int) {
	if len(c.Points) > 0 {
		if len(c.ConflictingContours) > 0 {
			pad := strings.Repeat(" ", offset)
			fmt.Printf("%s%d : len %d :  no kids %d : no col %d : col with parent %+v\n", pad, c.Id, len(c.Points), len(c.Children), len(c.ConflictingContours), c.ParentCollision)
		}
	}

	for _, ch := range c.Children {
		ContourStatsWithCollisions(ch, offset+2)
	}
}
