package border

import (
	"image"
	"strings"

	log "github.com/sirupsen/logrus"
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

	// usable or not. Not filtering out but marking that we may not use it. (say if we're conflicting with another contour)
	Usable bool
}

// NewContour create new contour
func NewContour(id int) *Contour {
	c := Contour{}
	c.Id = id
	c.BorderType = Hole
	c.ConflictingContours = make(map[int]bool)
	c.Usable = true
	return &c
}

// GetAllPoints returns all points in the contour and all children.
func (c *Contour) GetAllPoints() []image.Point {

	var allPoints []image.Point

	for _, p := range c.Points {
		allPoints = append(allPoints, p)
	}

	for _, ch := range c.Children {
		points := ch.GetAllPoints()
		allPoints = append(allPoints, points...)
	}

	return allPoints
}

// ContourStats generates writes the stats to a log about the contour and all children.
// Primarily used for debugging
func ContourStats(c *Contour, offset int) {
	if len(c.Points) > 0 {
		pad := strings.Repeat(" ", offset)
		log.Debugf("%s%d : len %d :  no kids %d : no col %d : col with parent %+v\n", pad, c.Id, len(c.Points), len(c.Children), len(c.ConflictingContours), c.ParentCollision)
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
			log.Debugf("%s%d : len %d :  no kids %d : no col %d : col with parent %+v\n", pad, c.Id, len(c.Points), len(c.Children), len(c.ConflictingContours), c.ParentCollision)
		}
	}

	for _, ch := range c.Children {
		ContourStatsWithCollisions(ch, offset+2)
	}
}
