package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
	"sort"
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
// fixed slice of 8, knocked down runtime of this by 40% or so.
func getValuesAroundPoint(borders *SuzukiImage, p image.Point) []int {

	pointVal := make([]int, 8, 8)
	count := 0
	minX := p.X - 1
	maxX := p.X + 2
	minY := p.Y - 1
	maxY := p.Y + 2
	for i := minY; i < maxY; i++ {
		for j := minX; j < maxX; j++ {

			// dont want centre.
			if !(i == p.Y && j == p.X) {
				pp := borders.GetXY(j, i)
				if pp != 0 {
					pp = 1
				}
				pointVal[count] = pp
				count++
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
		//fmt.Printf("XXXXXXXXXXX returning empty point!\n")
		//return image.Point{}, false
		return i2j2, false
	}

	result = (vv + cwRollDict[dir]) % 8
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
		//fmt.Printf("XXXXXXX returning empty point2\n")
		//return image.Point{}, false
		return i2j2, false
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

	result := (vv - ccwRollDict[dir] + 8) % 8
	p := ccwPixelDict[result]
	pp := centre.Sub(p)

	if pp.X == 0 && pp.Y == 0 {
		//fmt.Printf("XXXXX\n")
	}
	return pp, pixelFound
}

func findBorders(img *SuzukiImage) (*SuzukiImage, *Contours, int) {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	nbd := 1

	borders := img // reference to image?

	contours := NewContours()

	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			if borders.GetXY(j, i) != 0 {

				if borders.GetXY(j, i) == 1 && borders.GetXY(j-1, i) == 0 {
					nbd++
					i2j2 := image.Point{j - 1, i}
					i1j1, found := findClockwise(borders, image.Point{j, i}, i2j2)
					if found {
						if i == 1 && j == 4813 {
							//fmt.Printf("XXXX\n")
						}
						i2j2 = i1j1
						i3j3 := image.Point{j, i}
						count := 0
						for {
							count++
							i4j4, nextPixelFound := findCounterClockwise(borders, i3j3, i2j2)

							if nextPixelFound {
								borders.Set(i3j3, -1*nbd)
								contours.AddPointToContourId(nbd, i3j3)
							}
							if !nextPixelFound && borders.Get(i3j3) != 0 {
								borders.Set(i3j3, nbd)
								contours.AddPointToContourId(nbd, i3j3)
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
						contours.AddPointToContourId(nbd, image.Point{j, i})
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
								//fmt.Printf("22222\n")
								if nextPixelFound {
									borders.Set(i3j3, -1*nbd)
									contours.AddPointToContourId(nbd, i3j3)
								}
								if !nextPixelFound && borders.Get(i3j3) != 0 {
									borders.Set(i3j3, nbd)
									contours.AddPointToContourId(nbd, i3j3)
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
							contours.AddPointToContourId(nbd, image.Point{j, i})
						}
					}
				}
			}
		}
	}
	return borders, contours, nbd
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

	//black := color.RGBA{0, 0, 0, 255}

	padding := 2
	halfPadding := padding / 2

	padding = 0
	halfPadding = 0

	// need border to be black. Pad edges with 1 black pixel
	si := NewSuzukiImage(img.Bounds().Dx()+padding, img.Bounds().Dy()+padding)

	// dumb... but convert to own image format for now.
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			cc := 0
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			//fmt.Printf("%d %d %d %d\n", r, g, b, a)
			if !(r == 0 && g == 0 && b == 0) {
				cc = 1
			}
			si.SetXY(x+halfPadding, y+halfPadding, cc)
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

func saveContoursImage(filename string, c *Contours, width int, height int, flipBook bool, minContourSize int, smallestToLargest bool) error {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// naive fill
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.Black)
		}
	}

	_ = []color.RGBA{
		{50, 0, 0, 255},
		{100, 0, 0, 255},
		{150, 0, 0, 255},
		{200, 0, 0, 255},
		{250, 0, 0, 255},
		{50, 50, 0, 255},
		{100, 50, 0, 255},
		{150, 50, 0, 255},
		{200, 50, 0, 255},
		{250, 50, 0, 255},
		{50, 100, 0, 255},
		{100, 100, 0, 255},
		{150, 100, 0, 255},
		{200, 100, 0, 255},
		{250, 100, 0, 255},
		{50, 150, 0, 255},
		{100, 150, 0, 255},
		{150, 150, 0, 255},
		{200, 150, 0, 255},
		{250, 150, 0, 255},
		{50, 200, 0, 255},
		{100, 200, 0, 255},
		{150, 200, 0, 255},
		{200, 200, 0, 255},
		{250, 200, 0, 255},
	}

	colours := []color.RGBA{
		{255, 0, 0, 255},
		{255, 106, 0, 255},
		{255, 216, 0, 255},
		{0, 255, 0, 255},
		{127, 255, 197, 255},
		{72, 0, 255, 255},
		{255, 127, 182, 255},
	}

	max := len(colours)
	colour := 0
	count := 0

	contours := []*Contour{}
	for _, cc := range c.contours {

		contours = append(contours, cc)

		// only get length 11910
		if len(cc.points) == 11910 {
			//contours = append(contours, cc)
		}

	}

	//contours := c.contours
	if smallestToLargest {
		sort.Slice(contours, func(i int, j int) bool {
			return len(contours[i].points) < len(contours[j].points)
		})
	}

	for _, contour := range contours {

		fmt.Printf("contour %d has %d points\n", count, len(contour.points))
		if len(contour.points) < minContourSize {
			continue
		}
		colourToUse := colours[colour]

		for _, p := range contour.points {
			img.Set(p.X, p.Y, colourToUse)
		}
		colour++
		if colour >= max {
			colour = 0
		}

		// save new image per contour added...  crazy
		if flipBook {
			fn := fmt.Sprintf("%s-%d.png", filename, count)
			f, _ := os.Create(fn)
			png.Encode(f, img)
			f.Close()
		}
		count++
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)
	return nil
}

// writes details out to stdout
func displayContourStats(c *Contours) {

	shortestLength := 1000
	longestLength := 0
	averageLength := 0
	for _, cc := range c.contours {
		l := len(cc.points)
		if l > longestLength {
			longestLength = l
		}

		if l < shortestLength {
			shortestLength = l
		}
		averageLength += l
	}

	fmt.Printf("number of contours %d\n", len(c.contours))
	fmt.Printf("longest length %d\n", longestLength)
	fmt.Printf("shortest length %d\n", shortestLength)
	fmt.Printf("average length %d\n", averageLength/len(c.contours))
}

func main1() {
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
	si = loadImage("big-test-image.png")
	//si = loadImage("test5.png")
	//si = loadImage("small.png")
	//si = loadImage("image2.png")
	processingStart := time.Now()
	si2, contours, _ := findBorders(si)
	fmt.Printf("processing took %d ms\n", time.Now().Sub(processingStart).Milliseconds())

	displayContourStats(contours)
	//t := si2.DisplayAsText()
	//fmt.Printf("%+v\n", t)
	saveImage("border.png", si2)
	saveContoursImage("contours.png", contours, si2.Width, si2.Height, false, 30000, false)
	//saveContoursImage("./contours", contours, si2.Width, si2.Height, true, 0, false)
	fmt.Printf("load to save took %d ms\n", time.Now().Sub(start).Milliseconds())
}
