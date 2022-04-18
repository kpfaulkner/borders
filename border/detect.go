package border

import (
	"image"
)

var (

	// dirDelta determines which direction will we move based on the direction (0-7) index
	dirDelta = []image.Point{{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}}
)

// clockwise determines direction if we have 'dir' and turn clockwise
func clockwise(dir int) int {
	return (dir + 1) % 8
}

// counterClockwise determines direction if we have 'dir' and turn counterclockwise
func counterClockwise(dir int) int {
	return (dir + 7) % 8
}

// move moves the current point (pixel) in the direction 'dir'
func move(pixel image.Point, img *SuzukiImage, dir int) image.Point {
	newP := pixel.Add(dirDelta[dir])
	width := img.Width
	height := img.Height

	if (0 < newP.Y && newP.Y < height) && (0 < newP.X && newP.X < width) {
		if img.Get(newP) != 0 {
			return newP
		}
	}
	return image.Point{0, 0}
}

// calcDir returns index of dirDelta that matches direction taken.
func calcDir(from image.Point, to image.Point) int {
	delta := to.Sub(from)
	for i, d := range dirDelta {
		if d.X == delta.X && d.Y == delta.Y {
			return i
		}
	}

	// unsure... blow up for now.
	panic("BOOOOOOM cant figure out direction")
}

// createBorder returns the slice of Points making up the border/contour
// Also returns list of nbd's that are colliding with this. Can use to help create
// tree with collision info later.
func createBorder(img *SuzukiImage, p0 image.Point, p2 image.Point, nbd int, done []bool) ([]image.Point, map[int]bool) {

	// track which borders have conflicts
	collisionIndicies := make(map[int]bool)

	border := []image.Point{}
	dir := calcDir(p0, p2)
	moved := clockwise(dir)
	p1 := image.Point{0, 0}
	for moved != dir {
		newP := move(p0, img, moved)
		if newP.Y != 0 {
			p1 = newP
			break
		}
		moved = clockwise(moved)
	}

	if p1.X == 0 && p1.Y == 0 {
		return []image.Point{}, collisionIndicies
	}
	p2 = p1
	p3 := p0

	done = []bool{false, false, false, false, false, false, false, false}
	for {
		dir = calcDir(p3, p2)
		moved = counterClockwise(dir)
		p4 := image.Point{0, 0}
		done = []bool{false, false, false, false, false, false, false, false}
		for {
			p4 = move(p3, img, moved)
			if p4.Y != 0 {
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
			collisionIndicies[curP3] = true
			collisionIndicies[absNbd] = true
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

	return border, collisionIndicies
}

// addCollisionFlag mark contours with collisions with other contours.
func addCollisionFlag(contour *Contour, parentId int, contours map[int]*Contour, collisionIndices map[int]bool) {
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

// FindContours takes a SuzukiImage (basic 2d slice) and determines the Contours that are present.
// It returns the single parent contour which in turn has all other contours as children or further
// generations.
func FindContours(img *SuzukiImage) *Contour {

	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	nbd := 1
	lnbd := 1

	contours := make(map[int]*Contour)
	done := []bool{false, false, false, false, false, false, false, false}

	contour := NewContour(1)
	contours[lnbd] = contour

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
				border, collectionIndices := createBorder(img, p0, from, nbd, done)
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
				contours[nbd] = contour
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
	return contours[1]
}
