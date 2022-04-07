package border

import "image"

const (
	Outer = 0
	Hole  = 1
)

type Contour struct {
	Points     []image.Point
	Id         int
	BorderType int
	ParentId   int
}

func NewContour(id int) *Contour {
	c := Contour{}
	c.Id = id
	//c.Points = make([]image.Point, 500, 1000)
	c.BorderType = Hole
	return &c
}

func (c *Contour) AddPoint(p image.Point) error {
	c.Points = append(c.Points, p)
	return nil
}

type Contours struct {
	contours map[int]*Contour
}

func NewContours() *Contours {
	c := Contours{}
	c.contours = make(map[int]*Contour, 100)
	return &c
}

func (c *Contours) AddPointToContourId(id int, p image.Point) error {
	if _, has := c.contours[id]; !has {
		c.contours[id] = NewContour(id)
	}
	c.contours[id].AddPoint(p)
	return nil
}
