package border

import (
	"errors"
	"image"

	"github.com/kpfaulkner/borders/common"
)

// dirDelta determines which direction we move based on the direction (0-7) index.
// Fixed-size array so the type itself documents the 8-direction invariant and the
// compiler can prove dirDelta[dir] is in bounds for any `dir` it can range-check.
var dirDelta = [8]image.Point{{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}}

// FindContours takes a SuzukiImage and determines the Contours that are present.
// It returns the single parent contour which in turn has all other contours as children or further
// generations.
func FindContours(img *common.SuzukiImage) (*Contour, error) {
	nbd := 1
	lnbd := 1

	// contours[id] holds the contour with that id. Contour ids are assigned
	// sequentially starting at 1, so direct slice indexing replaces a map.
	// Index 0 is a nil sentinel — no contour ever has id 0.
	contours := []*Contour{nil, NewContour(1)}

	height := img.Height
	width := img.Width

	for i := 0; i < height; i++ {
		lnbd = 1
		for j := 0; j < width; j++ {
			fji := img.GetXY(j, i)
			isOuter := fji == 1 && (j == 0 || img.GetXY(j-1, i) == 0)
			isHole := fji >= 1 && (j == width-1 || img.GetXY(j+1, i) == 0)
			if isOuter || isHole {

				var contourPrime *Contour
				contour := NewContour(1)
				from := image.Point{j, i}
				parentId := 0
				if isOuter {
					nbd += 1
					from = from.Sub(image.Point{1, 0})
					contour.BorderType = Outer
					contourPrime = contours[lnbd]
					if contourPrime.BorderType == Outer {
						parentId = contourPrime.ParentId
					} else {
						parentId = contourPrime.Id
					}
				} else {
					nbd += 1
					if fji > 1 {
						lnbd = fji
					}
					contourPrime = contours[lnbd]
					from = from.Add(image.Point{1, 0})
					contour.BorderType = Hole
					if contourPrime.BorderType == Outer {
						parentId = contourPrime.Id
					} else {
						parentId = contourPrime.ParentId
					}
				}

				p0 := image.Point{j, i}
				border, collectionIndices, err := createBorder(img, p0, from, nbd)
				if err != nil {
					return nil, err
				}

				if len(border) == 0 {
					border = append(border, p0)
					img.Set(p0, -1*nbd)
				}

				if parentId != 0 {
					parent := contours[parentId]
					parent.Children = append(parent.Children, contour)
					contour.Parent = contours[parentId]
				}
				contour.ParentId = parentId
				contour.Points = border
				contour.Id = nbd
				// nbd was incremented to len(contours), so append lands the new contour at index nbd.
				contours = append(contours, contour)
				addCollisionFlag(contour, parentId, contours, collectionIndices)
			}
			if fji != 0 && fji != 1 {
				lnbd = fji
				if lnbd < 0 {
					lnbd *= -1
				}
			}
		}
	}

	finalContour := contours[1]

	// image was padded... so now shift every co-ord by -1,-1
	if img.HasPadding() {
		shiftContour(finalContour)
	}
	return finalContour, nil
}

// clockwise determines direction if we have 'dir' and turn clockwise
func clockwise(dir int) int {
	return (dir + 1) % 8
}

// counterClockwise determines direction if we have 'dir' and turn counterclockwise
func counterClockwise(dir int) int {
	return (dir + 7) % 8
}

// move moves the current point (pixel) in the direction 'dir'.
// Returns ok=false when the neighbour is out of bounds or background.
func move(pixel image.Point, img *common.SuzukiImage, dir int) (image.Point, bool) {
	newP := pixel.Add(dirDelta[dir])
	if newP.X < 0 || newP.X >= img.Width || newP.Y < 0 || newP.Y >= img.Height {
		return image.Point{}, false
	}
	if img.Get(newP) == 0 {
		return image.Point{}, false
	}
	return newP, true
}

// calcDir returns index of dirDelta that matches direction taken.
func calcDir(from image.Point, to image.Point) (int, error) {
	delta := to.Sub(from)
	for i, d := range dirDelta {
		if d.X == delta.X && d.Y == delta.Y {
			return i, nil
		}
	}

	return 0, errors.New("unable to determine direction")
}

// createBorder returns the slice of Points making up the border/contour
// Also returns list of nbd's that are colliding with this. Can use to help create
// tree with collision info later.
func createBorder(img *common.SuzukiImage, p0 image.Point, p2 image.Point, nbd int) ([]image.Point, map[int]bool, error) {

	// track which borders have conflicts
	collisionIndices := make(map[int]bool)

	border := []image.Point{}
	dir, err := calcDir(p0, p2)
	if err != nil {
		return nil, nil, err
	}

	moved := clockwise(dir)
	var p1 image.Point
	foundP1 := false
	for moved != dir {
		if newP, ok := move(p0, img, moved); ok {
			p1 = newP
			foundP1 = true
			break
		}
		moved = clockwise(moved)
	}

	if !foundP1 {
		return []image.Point{}, collisionIndices, nil
	}
	p2 = p1
	p3 := p0

	for {
		dir, err = calcDir(p3, p2)
		if err != nil {
			return nil, nil, err
		}
		moved = counterClockwise(dir)
		var p4 image.Point
		var done [8]bool
		for {
			var ok bool
			p4, ok = move(p3, img, moved)
			if ok {
				break
			}
			done[moved] = true
			moved = counterClockwise(moved)
		}

		// detect if colliding with something else (ie not 0 nor 1)
		curP3 := img.Get(p3)
		if curP3 != 1 {
			if curP3 < 0 {
				curP3 *= -1
			}

			absNbd := nbd
			if absNbd < 0 {
				absNbd *= -1
			}
			collisionIndices[curP3] = true
			collisionIndices[absNbd] = true
		}

		border = append(border, p3)
		if p3.Y == img.Height-1 || done[2] {
			img.Set(p3, -1*nbd)
		} else if img.Get(p3) == 1 {
			img.Set(p3, nbd)
		}

		if p4.X == p0.X && p4.Y == p0.Y && p3.X == p1.X && p3.Y == p1.Y {
			break
		}

		p2 = p3
		p3 = p4
	}

	return border, collisionIndices, nil
}

// addCollisionFlag mark contours with collisions with other contours.
func addCollisionFlag(contour *Contour, parentId int, contours []*Contour, collisionIndices map[int]bool) {
	for contour1 := range collisionIndices {

		// quick indicator to say colliding with parent.
		if contour1 == parentId {
			contour.ParentCollision = true
		}

		for contour2 := range collisionIndices {
			if contour1 != contour2 {
				c1 := contours[contour1]
				c1.ConflictingContours[contour2] = true
			}
		}
	}
}

func shiftContour(contour *Contour) {
	for i, _ := range contour.Points {
		contour.Points[i].X = contour.Points[i].X - 1
		contour.Points[i].Y = contour.Points[i].Y - 1
	}

	for _, child := range contour.Children {
		shiftContour(child)
	}
}
