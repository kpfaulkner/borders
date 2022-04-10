package border

import "image"

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
