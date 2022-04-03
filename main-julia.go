package main

import (
	"fmt"
	"image"
	"time"

	"github.com/pkg/profile"
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

// returns index of dirDelta that matches direction taken.
func fromTo(from image.Point, to image.Point) int {
	delta := to.Sub(from)
	for i, d := range dirDelta {
		if d.X == delta.X && d.Y == delta.Y {
			return i
		}
	}

	// unsure... blow up.
	panic("BOOOOOOM cant figure out direction")
}

func detectMove(img *SuzukiImage, p0 image.Point, p2 image.Point, nbd int, done []bool) []image.Point {
	border := []image.Point{}
	dir := fromTo(p0, p2)
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
		dir = fromTo(p3, p2)
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

func findContours(img *SuzukiImage) []*Contour {

	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	nbd := 1
	contours := []*Contour{}
	done := []bool{false, false, false, false, false, false, false, false}

	height := img.Height
	width := img.Width
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			isOuter := img.GetXY(j, i) == 1 && (j == 0 || img.GetXY(j-1, i) == 0)
			isHole := img.GetXY(j, i) >= 1 && (j == width-1 || img.GetXY(j+1, i) == 0)
			if isOuter || isHole {

				from := image.Point{j, i}
				if isOuter {
					nbd++
					from = from.Sub(image.Point{1, 0})
				} else {
					nbd++
					from = from.Add(image.Point{1, 0})
				}

				p0 := image.Point{j, i}
				border := detectMove(img, p0, from, nbd, done)
				if len(border) == 0 {
					border = append(border, p0)
					img.Set(p0, -1*nbd)
				}
				contour := NewContour(nbd)
				contour.points = border
				contours = append(contours, contour)
			}
		}
	}
	return contours
}

func main() {
	fmt.Printf("So it begins...\n")

	//img := loadImage("image2.png")
	img := loadImage("big-test-image.png")

	start := time.Now()
	cont := findContours(img)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	saveContourSliceImage("julia-contour.png", cont, img.Width, img.Height, false, 0, false)

	fmt.Printf("NUm contours are %d\n", len(cont))
}
