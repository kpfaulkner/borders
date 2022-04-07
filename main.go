package main

import (
	"fmt"
	"image"
	"sort"
	"time"
)

var (
	dirDelta = []image.Point{{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}}
)

func clockwise(dir int) int {
	return ((dir + 1) % 8)
}

func counterClockwise(dir int) int {
	return (dir + 7) % 8
}

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

	// unsure... blow up.
	panic("BOOOOOOM cant figure out direction")
}

// createBorder returns the slice of points making up the border/contour
func createBorder(img *SuzukiImage, p0 image.Point, p2 image.Point, nbd int, done []bool) []image.Point {
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
		return []image.Point{}
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
	return border
}

func findContours(img *SuzukiImage) map[int]*Contour {

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

			isOuter := img.GetXY(j, i) == 1 && (j == 0 || img.GetXY(j-1, i) == 0)
			isHole := img.GetXY(j, i) >= 1 && (j == width-1 || img.GetXY(j+1, i) == 0)
			if isOuter || isHole {

				var contourPrime *Contour
				contour := NewContour(1)
				from := image.Point{j, i}
				if isOuter {
					nbd += 1
					from = from.Sub(image.Point{1, 0})
					contour.borderType = Outer
					contourPrime = contours[lnbd]
					if contourPrime.borderType == Outer {
						contour.parentId = contourPrime.parentId
					} else {
						contour.parentId = contourPrime.id
					}
				} else {
					nbd += 1
					if fji > 1 {
						lnbd = fji
					}
					contourPrime = contours[lnbd]
					from = from.Add(image.Point{1, 0})
					contour.borderType = Hole
					if contourPrime.borderType == Outer {
						contour.parentId = contourPrime.id
					} else {
						contour.parentId = contourPrime.parentId
					}
				}

				p0 := image.Point{j, i}
				border := createBorder(img, p0, from, nbd, done)
				if len(border) == 0 {
					border = append(border, p0)
					img.Set(p0, -1*nbd)
				}

				contour.points = border
				contour.id = nbd
				contours[nbd] = contour
			}
			if fji != 0 && fji != 1 {
				lnbd = fji
				if lnbd < 0 {
					lnbd *= -1
				}
			}
		}
	}
	return contours
}

func main() {
	fmt.Printf("So it begins...\n")

	//img := loadImage("image1.png")
	img := loadImage("big-test-image.png")

	start := time.Now()
	cont := findContours(img)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	//saveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0, false)
	saveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0, false)
	//saveContourSliceImage("c:/temp/contour/contour", cont, img.Width, img.Height, true, 0, false)

	contours := []*Contour{}
	for _, cc := range cont {
		contours = append(contours, cc)
	}

	sort.Slice(contours, func(i int, j int) bool {
		return contours[i].id < contours[j].id
	})

	for _, c := range contours {
		parentId := c.parentId

		if c.borderType == Outer {
			if c.parentId != 1 {
				parent := cont[c.parentId]
				if parent.borderType == Outer {
					// how is parent of an outer... and outer? switch to parent being 2?
					parentId = 1
				}
			}
		}
		fmt.Printf("%d %d : %d\n", c.id, parentId, c.borderType)
	}

	fmt.Printf("Num contours are %d\n", len(cont))
}
