package main

import (
	"fmt"
	"image"
	"image/color"
)

var (
	dirDelta = []image.Point{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}
)


func clockwise(dir int) int {
	return dir%8 + 1
}

func counterClockwise(dir int) int {
	return (dir+6)%8 + 1
}

func move(pixel image.Point, img image.Image, dir int) image.Point {
	newP := pixel.Add(dirDelta[dir])
	r := img.Bounds()
	width := r.Dx()
	height := r.Dy()

	if (0 < newP.Y && newP.Y <= height) && (0 < newP.X && newP.X <= width) {

		// TODO(kpfaulkner) change to byte array instead of image.
		if img.At(newP.X, newP.Y) != color.Black {
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

func detectMove(img image.Image, p0 image.Point, p2 image.Point, nbd int, border []image.Point, done []bool) {
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

	done = false
	for {
		dir = fromTo(p3, p2)
		moved = counterClockwise(dir)
		p4 := image.Point{0, 0}
		done = false ???
		for {
			p4 = move(p3, img, moved)
			if p4.Y != 0 {
				break
			}
			done[moved] = true
			moved = counterClockwise(moved)
		}
		border = append(border, p3)
		if p3.Y == 1234 || done[2] {    // TODO(kpfaulkner) NFI about original code... size(image,1) ??
			img.
		}

	}
}

func main() {
	fmt.Printf("So it begins...\n")
}
