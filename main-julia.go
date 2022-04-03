package main

import (
	"fmt"
	"image"
	"math"
	"time"
)

var (
	//dirDelta = []image.Point{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}
	dirDelta    = []image.Point{{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}}
	fromToCount = 0
)

func clockwise(dir int) int {
	//return (dir % 7) + 1
	return ((dir + 1) % 8)
}

func counterClockwise(dir int) int {
	//return (dir + 6) % 7
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
	//fmt.Printf("from to count %d\n", fromToCount)
	fromToCount++
	//fmt.Printf("from %+v : to %+v\n", from, to)
	delta := to.Sub(from)
	//fmt.Printf("delta %+v\n", delta)
	for i, d := range dirDelta {
		if d.X == delta.X && d.Y == delta.Y {
			return i
		}
	}

	// unsure... blow up.
	panic("BOOOOOOM cant figure out direction")
}

func detectMove(img *SuzukiImage, p0 image.Point, p2 image.Point, nbd int, done []bool) []image.Point {
	//fmt.Printf("detectMove\n")
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

	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	nbd := 1
	contourList := [][]image.Point{}
	contours := []*Contour{}
	done := []bool{false, false, false, false, false, false, false, false}

	height := img.Height
	width := img.Width
	var lnbd int
	for i := 0; i < height; i++ {
		lnbd = 1
		for j := 0; j < width; j++ {
			fji := img.GetXY(j, i)
			isOuter := img.GetXY(j, i) == 1 && (j == 0 || img.GetXY(j-1, i) == 0)
			isHole := img.GetXY(j, i) >= 1 && (j == width-1 || img.GetXY(j+1, i) == 0)
			if isOuter || isHole {

				// fmt.Printf("outer %+v : hole %+v\n", isOuter, isHole)
				//border := []image.Point{}
				from := image.Point{j, i}
				//fmt.Printf("FROM is %+v\n", from)
				if isOuter {
					nbd++
					from = from.Sub(image.Point{1, 0})
				} else {
					nbd++
					if fji > 1 {
						lnbd = fji
					}
					from = from.Add(image.Point{1, 0})
				}

				p0 := image.Point{j, i}
				border := detectMove(img, p0, from, nbd, done)
				if len(border) == 0 {
					border = append(border, p0)
					img.Set(p0, -1*nbd)
				}
				//fmt.Printf("borderyyy length %d\n", len(border))
				contour := NewContour(nbd)
				contour.points = border
				contours = append(contours, contour)
				contourList = append(contourList, border)
				//fmt.Printf("contour length %d\n", len(contourList))
			}
			if fji != 0 && fji != 1 {
				lnbd = int(math.Abs(float64(fji)))
			}
		}
	}

	fmt.Printf("lnbd %d\n", lnbd)
	return contours
}

func main() {
	fmt.Printf("So it begins...\n")

	/*
		img := NewSuzukiImage(100, 100)
		for x := 40; x < 60; x++ {
			for y := 40; y < 60; y++ {
				img.SetXY(x, y, 1)
			}
		} */

	//img := loadImage("image2.png")
	img := loadImage("big-test-image.png")

	fmt.Printf("A11 %d\n", img.GetXY(4811, 0))
	fmt.Printf("A12 %d\n", img.GetXY(4812, 0))
	fmt.Printf("A13 %d\n", img.GetXY(4813, 0))
	fmt.Printf("A14 %d\n", img.GetXY(4814, 0))
	fmt.Printf("A15 %d\n", img.GetXY(4815, 0))
	fmt.Printf("A16 %d\n", img.GetXY(4816, 0))
	fmt.Printf("B11 %d\n", img.GetXY(4811, 1))
	fmt.Printf("B12 %d\n", img.GetXY(4812, 1))
	fmt.Printf("B13 %d\n", img.GetXY(4813, 1))
	fmt.Printf("B14 %d\n", img.GetXY(4814, 1))
	fmt.Printf("B15 %d\n", img.GetXY(4815, 1))
	fmt.Printf("B16 %d\n", img.GetXY(4816, 1))

	start := time.Now()
	cont := findContours(img)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	saveContourSliceImage("julia-contour.png", cont, img.Width, img.Height, false, 0, false)

	fmt.Printf("NUm contours are %d\n", len(cont))
}
