package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
	"sort"
)

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

	padding = 1
	halfPadding = 1

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

func saveContourSliceImage(filename string, c map[int]*Contour, width int, height int, flipBook bool, minContourSize int, smallestToLargest bool) error {

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
	for _, cc := range c {
		contours = append(contours, cc)
	}

	sort.Slice(contours, func(i int, j int) bool {
		return contours[i].id < contours[j].id
	})

	//contours := c.contours
	if smallestToLargest {
		sort.Slice(contours, func(i int, j int) bool {
			return len(contours[i].points) < len(contours[j].points)
		})
	}

	for _, contour := range contours {
		fmt.Printf("contour %d has %d points : borderType %d\n", count, len(contour.points), contour.borderType)
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
