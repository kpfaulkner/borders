package main

import (
	"fmt"
	"image"
)

func findClockwise(borders *SuzukiImage, centre []image.Point,, i2j2 []image.Point) (image.Point, bool) {
	mask := [][]bool{{true, true, true}, {true, false, true}, {true, true, true}}
	_ = mask

	return image.Point{}, true
}

func findCounterClockwise(borders *SuzukiImage, centre []image.Point,, i2j2 []image.Point) (image.Point, bool) {

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
					i1j1, found := findClockwise(borders, []image.Point{{j, i}}, []image.Point{i2j2})
					if found {
						i2j2 = i1j1
						i3j3 := image.Point{j, i}
						for {
							i4j4, nextPixelFound := findCounterClockwise(borders, []image.Point{{j, i}}, []image.Point{i2j2})
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
						i1j1, found := findClockwise(borders, []image.Point{{j, i}}, []image.Point{i2j2})
						if found {
							i2j2 = i1j1
							i3j3 := image.Point{j, i}
							for {
								i4j4, nextPixelFound := findCounterClockwise(borders, []image.Point{i3j3}, []image.Point{i2j2})
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

}
