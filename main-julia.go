package main

import (
	"fmt"
	"image"
	"math"
)

var (
	//dirDelta = []image.Point{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}
	dirDelta = []image.Point{{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}}
)

func clockwise(dir int) int {
	return dir%7 + 1
}

func counterClockwise(dir int) int {
	return (dir+5)%7 + 1
}

func move(pixel image.Point, img *SuzukiImage, dir int) image.Point {
	newP := pixel.Add(dirDelta[dir])
	width := img.Width
	height := img.Height

	if (0 < newP.Y && newP.Y <= height) && (0 < newP.X && newP.X <= width) {
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

func detectMove(img *SuzukiImage, p0 image.Point, p2 image.Point, nbd int, border []image.Point, done []bool) {
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
		return
	}
	p2 = p1
	p3 := p0

	// same as julia done .= false   I think
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
		if p3.Y == 1234 || done[2] { // TODO(kpfaulkner) NFI about original code... size(image,1) ??
			img.Set(p3, -1*nbd)
		} else if img.Get(p3) == 1 {
			img.Set(p3, nbd)
		}

		if p4 == p0 && p3 == p1 {
			break
		}

		p2 = p3
		p3 = p4
	}
}

func findContours(img *SuzukiImage) []image.Point {
	nbd := 1
	contourList := []image.Point{}
	done := []bool{false, false, false, false, false, false, false, false}

	height := img.Height
	width := img.Width
	var lnbd int
	for i := 0; i < height; i++ {
		lnbd = 1
		for j := 0; j < width; j++ {
			fji := img.GetXY(j, i)
			isOuter := img.GetXY(j, i) == 1 && (j == 0 || img.GetXY(j-1, i) == 0)
			isHole := img.GetXY(j, i) >= 1 && (j == width || img.GetXY(j+1, i) == 0)
			if isOuter || isHole {
				border := []image.Point{}
				from := image.Point{j, i}
				if isOuter {
					nbd++
					from = from.Sub(image.Point{1, 0})
				} else {
					nbd++
					if fji > 1 {
						lnbd = fji
					}
					from.Add(image.Point{1, 0})
				}

				p0 := image.Point{j, i}
				detectMove(img, p0, from, nbd, border, done)
				if len(border) == 0 {
					border = append(border, p0)
					img.Set(p0, -1*nbd)
				}
				contourList = append(contourList, border...) // TODO(kpfaulkner) check this!
			}
			if fji != 0 && fji != 1 {
				lnbd = int(math.Abs(float64(fji)))
			}
		}
	}

	fmt.Printf("LNBD is %d\n", lnbd)
	return contourList
}

func mainjulia() {
	fmt.Printf("So it begins...\n")

	img := NewSuzukiImage(100, 100)
	for x := 40; x < 60; x++ {
		for y := 40; y < 60; y++ {
			img.SetXY(x, y, 1)
		}
	}

	cont := findContours(img)
	fmt.Printf("Contours are %+v\n", cont)
}
