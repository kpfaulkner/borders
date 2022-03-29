package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
	"time"
)

var (
	cwRollDict  = make(map[image.Point]int)
	cwPixelDict = make([]image.Point, 8)

	ccwRollDict  = make(map[image.Point]int)
	ccwPixelDict = make([]image.Point, 8)
)

func init() {
	cwRollDict[image.Point{1, 1}] = 0
	cwRollDict[image.Point{0, 1}] = 1
	cwRollDict[image.Point{-1, 1}] = 2
	cwRollDict[image.Point{-1, 0}] = 3
	cwRollDict[image.Point{-1, -1}] = 4
	cwRollDict[image.Point{0, -1}] = 5
	cwRollDict[image.Point{1, -1}] = 6
	cwRollDict[image.Point{1, 0}] = 7

	cwPixelDict[0] = image.Point{1, 1}
	cwPixelDict[1] = image.Point{0, 1}
	cwPixelDict[2] = image.Point{-1, 1}
	cwPixelDict[3] = image.Point{-1, 0}
	cwPixelDict[4] = image.Point{-1, -1}
	cwPixelDict[5] = image.Point{0, -1}
	cwPixelDict[6] = image.Point{1, -1}
	cwPixelDict[7] = image.Point{1, 0}

	ccwRollDict[image.Point{1, 1}] = 0
	ccwRollDict[image.Point{0, 1}] = 1
	ccwRollDict[image.Point{-1, 1}] = 2
	ccwRollDict[image.Point{-1, 0}] = 3
	ccwRollDict[image.Point{-1, -1}] = 4
	ccwRollDict[image.Point{0, -1}] = 5
	ccwRollDict[image.Point{1, -1}] = 6
	ccwRollDict[image.Point{1, 0}] = 7

	ccwPixelDict[0] = image.Point{1, 1}
	ccwPixelDict[1] = image.Point{1, 0}
	ccwPixelDict[2] = image.Point{1, -1}
	ccwPixelDict[3] = image.Point{0, -1}
	ccwPixelDict[4] = image.Point{-1, -1}
	ccwPixelDict[5] = image.Point{-1, 0}
	ccwPixelDict[6] = image.Point{-1, 1}
	ccwPixelDict[7] = image.Point{0, 1}

}

func rotateSlice(s []int, rotation int) []int {
	//rotation := v % len(s)
	if rotation < 0 {
		r := rotation * -1
		newS := append(s[r:], s[:r]...)
		return newS
	}

	index := len(s) - rotation
	newS := append(s[index:], s[:index]...)

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
				if pp != 0 {
					pp = 1
				}
				pointVal = append(pointVal, pp)
			}
		}
	}

	return pointVal

}

// naive approach but will do until perf testing says otherwise
func firstIndexContaining(l []int, v int) (int, error) {
	for i, vv := range l {
		if vv == v {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no value")
}

// steps:
// 1) get 3x3 grid with centre being the centre of the grid
// 2) swap... (unsure reason)
// 3) rotate
// 4)
func findClockwise(borders *SuzukiImage, centre image.Point, i2j2 image.Point) (image.Point, bool) {

	values := getValuesAroundPoint(borders, centre)

	// this is purely taken from existing code...  do NOT understand why yet!
	values[7], values[3], values[6], values[5], values[4] = values[3], values[4], values[5], values[6], values[7]

	dir := centre.Sub(i2j2)
	v := cwRollDict[dir]
	values2 := rotateSlice(values, -1*v)

	// anything non 0 set to 1
	for i, v := range values2 {
		if v != 0 {
			values2[i] = 1
		}
	}

	var result int
	dir = centre.Sub(i2j2)
	vv, err := firstIndexContaining(values2, 1)
	if err != nil {
		return image.Point{}, false
	}
	if vv+cwRollDict[dir] >= 8 {
		result = vv - 8 + cwRollDict[dir]
	} else {
		result = vv + cwRollDict[dir]
	}

	p := cwPixelDict[result]

	pp := centre.Sub(p)
	return pp, true
}

func findCounterClockwise(borders *SuzukiImage, centre image.Point, i2j2 image.Point) (image.Point, bool) {

	values := getValuesAroundPoint(borders, centre)

	// this is purely taken from existing code...  do NOT understand why yet!
	values[7], values[6], values[1], values[5], values[2], values[3], values[4] = values[1], values[2], values[3], values[4], values[5], values[6], values[7]
	dir := centre.Sub(i2j2)
	v := ccwRollDict[dir]
	values2 := rotateSlice(values, v)
	values2[0] = 0

	// anything non 0 set to 1
	for i, v := range values2 {
		if v != 0 {
			values2[i] = 1
		}
	}

	pixelFound := false
	dir = centre.Sub(i2j2)

	vv, err := firstIndexContaining(values2, 1)
	if err != nil {
		return image.Point{}, false
	}
	if ccwRollDict[dir] > 3 {
		if ccwRollDict[dir]-vv < 3 {
			pixelFound = true
		}
	}
	if ccwRollDict[dir] < 3 {
		if 8+ccwRollDict[dir]-vv < 3 {
			pixelFound = true
		}
	}

	var result int
	if vv-ccwRollDict[dir] < 0 {
		result = vv + 8 - ccwRollDict[dir]
	} else {
		result = vv - ccwRollDict[dir]
	}

	p := ccwPixelDict[result]
	pp := centre.Sub(p)
	return pp, pixelFound
}

func findBorders(img *SuzukiImage) (*SuzukiImage, int) {
	nbd := 1

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

func loadImage(filename string) *SuzukiImage {
	f, err := os.Open(filename)
	if err != nil {
		panic("BOOM on file")
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic("BOOM2 on file")
	}

	black := color.RGBA{0, 0, 0, 255}
	si := NewSuzukiImage(img.Bounds().Dx(), img.Bounds().Dy())
	// dumb... but convert to own image format for now.
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			cc := 0
			c := img.At(x, y)
			if c != black {
				cc = 1
			}
			si.SetXY(x, y, cc)
		}

	}

	return si
}

func saveImage(filename string, si *SuzukiImage) error {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{si.Width, si.Height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < si.Width; x++ {
		for y := 0; y < si.Height; y++ {
			p := si.GetXY(x, y)
			if p != 0 && p != 1 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)
	return nil
}

func main() {
	fmt.Printf("So it begins...\n")

	si := NewSuzukiImage(50, 50)

	for x := 10; x < 40; x++ {
		for y := 10; y < 40; y++ {
			si.SetXY(x, y, 1)
		}
	}

	for x := 42; x < 46; x++ {
		for y := 42; y < 46; y++ {
			si.SetXY(x, y, 1)
		}
	}

	// add hole 1
	for x := 20; x < 24; x++ {
		for y := 20; y < 24; y++ {
			si.SetXY(x, y, 0)
		}
	}

	// add hole 2
	for x := 28; x < 32; x++ {
		for y := 28; y < 32; y++ {
			si.SetXY(x, y, 0)
		}
	}

	start := time.Now()
	//si = loadImage("image2.png")
	processingStart := time.Now()
	si2, _ := findBorders(si)
	fmt.Printf("processing took %d ms\n", time.Now().Sub(processingStart).Milliseconds())
	t := si2.DisplayAsText()
	fmt.Printf("%+v\n", t)
	saveImage("border.png", si2)
	fmt.Printf("load to save took %d ms\n", time.Now().Sub(start).Milliseconds())
}
