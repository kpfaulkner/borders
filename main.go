package main

import (
	"fmt"
	"image"
)

var (
	rollDict  = make(map[image.Point]int)
	pixelDict = make([]image.Point, 8)
)

func init() {
	rollDict[image.Point{1, 1}] = 0
	rollDict[image.Point{0, 1}] = 1
	rollDict[image.Point{-1, 1}] = 2
	rollDict[image.Point{-1, 0}] = 3
	rollDict[image.Point{-1, -1}] = 4
	rollDict[image.Point{0, -1}] = 5
	rollDict[image.Point{1, -1}] = 6
	rollDict[image.Point{1, 0}] = 7

	pixelDict[0] = image.Point{1, 1}
	pixelDict[1] = image.Point{0, 1}
	pixelDict[2] = image.Point{-1, 1}
	pixelDict[3] = image.Point{-1, 0}
	pixelDict[4] = image.Point{-1, -1}
	pixelDict[5] = image.Point{0, -1}
	pixelDict[6] = image.Point{1, -1}
	pixelDict[7] = image.Point{1, 0}
}

func rotateSliceLeft(s []int, v int) []int {
	rotation := v % len(s)
	newS := append(s[rotation:], s[:rotation]...)
	return newS
}

// gets values around a point.
// filters out centre (p) point... so slice should be 8 elements in length.
func getValuesAroundPoint(borders *SuzukiImage, p image.Point) []int {

	pointVal := []int{}
	for i := p.Y - 1; i < p.Y+2; i++ {
		for j := p.X - 1; j < p.X+2; j++ {

			// dont want centre.
			if !(i == p.Y && j == p.X) {
				pp := borders.GetXY(j, i)
				if pp > 0 {
					pp = 1
				}
				pointVal = append(pointVal, pp)
			}
		}
	}

	return pointVal

}

// steps:
// 1) get 3x3 grid with centre being the centre of the grid
// 2) swap... (unsure reason)
// 3) rotate
// 4)
func findClockwise(borders *SuzukiImage, centre image.Point, i2j2 image.Point) (image.Point, bool) {

	values := getValuesAroundPoint(borders, centre)

	values[3] = 1
	values[5] = 1
	values[6] = 1
	values[7] = 1

	// this is purely taken from existing code...  do NOT understand why yet!
	values[7], values[3], values[6], values[5], values[4] = values[3], values[4], values[5], values[6], values[7]

	dir := centre.Sub(i2j2)
	v := rollDict[dir]
	values2 := rotateSliceLeft(values, v)

	var result int
	dir = centre.Sub(i2j2)
	if values2[1]+rollDict[dir] >= 8 {
		result = values2[1] - 8 + rollDict[dir]
	} else {
		result = values2[1] + rollDict[dir]
	}

	p := pixelDict[result]

	pp := centre.Sub(p)
	return pp, true
}

func findCounterClockwise(borders *SuzukiImage, centre image.Point, i2j2 image.Point) (image.Point, bool) {

	return image.Point{}, false
}

func findBorders(img *SuzukiImage) (*SuzukiImage, int) {
	nbd := 1

	// borders[borders == 255] = 1  NFI!!!
	borders := img // reference to image?

	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			if borders.GetXY(j, i) != 0 {

				if borders.GetXY(j, i) == 1 && borders.GetXY(j-1, i) == 0 {
					nbd++
					i2j2 := image.Point{j - 1, i}
					i1j1, found := findClockwise(borders, image.Point{j, i}, i2j2)
					if found {
						i2j2 = i1j1
						i3j3 := image.Point{j, i}
						for {
							i4j4, nextPixelFound := findCounterClockwise(borders, image.Point{j, i}, i2j2)
							if nextPixelFound {
								borders.Set(i3j3, -1*nbd)
							}
							if !nextPixelFound && borders.Get(i3j3) == 1 {
								borders.Set(i3j3, nbd)
							}

							if i4j4.X == j && i4j4.Y == i && i3j3.X == i1j1.X && i3j3.Y == i1j1.Y {
								break
							} else {
								i2j2 = i3j3
								i3j3 = i4j4
							}
						}
					} else {
						borders.SetXY(j, i, -1*nbd)
					}

				} else {
					if borders.GetXY(j, i) >= 1 && borders.GetXY(j+1, i) == 0 {
						nbd++
						i2j2 := image.Point{j + 1, i}
						i1j1, found := findClockwise(borders, image.Point{j, i}, i2j2)
						if found {
							i2j2 = i1j1
							i3j3 := image.Point{j, i}
							for {
								i4j4, nextPixelFound := findCounterClockwise(borders, i3j3, i2j2)
								if nextPixelFound {
									borders.Set(i3j3, -1*nbd)
								}
								if !nextPixelFound && borders.Get(i3j3) == 1 {
									borders.Set(i3j3, nbd)
								}

								if i4j4.X == j && i4j4.Y == i && i3j3.X == i1j1.X && i3j3.Y == i1j1.Y {
									break
								} else {
									i2j2 = i3j3
									i3j3 = i4j4
								}
							}
						} else {
							borders.SetXY(j, i, -1*nbd)
						}
					}
				}
			}
		}
	}
	return borders, nbd
}

func main() {
	fmt.Printf("So it begins...\n")

	si := NewSuzukiImage(100, 100)
	findClockwise(si, image.Point{292, 74}, image.Point{293, 74})
}
