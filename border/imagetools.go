package border

import (
	"fmt"
	"github.com/kpfaulkner/borders/common"
	image2 "github.com/kpfaulkner/borders/image"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
)

// LoadImage loads a PNG and returns a SuzukiImage.
// This may change since SuzukiImage may not really be required.
// erode flag forces the eroding of the image before converting to a SuzukiImage.
// This is to remove any "spikes" that may appear in the generated boundary.
func LoadImage(filename string, erode int, dilate int) (*common.SuzukiImage, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	// need border to be black. Pad edges with 1 black pixel
	si := common.NewSuzukiImage(img.Bounds().Dx(), img.Bounds().Dy())

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
			si.SetXY(x, y, cc)
		}

	}

	if erode != 0 {
		si, err = image2.Erode(si, erode)
		if err != nil {
			return nil, err
		}
	}

	if dilate != 0 {
		si, err = image2.Dilate(si, dilate)
		if err != nil {
			return nil, err
		}
	}

	return si, nil
}

// SaveImage saves a SuzukiImage to filename
func SaveImage(filename string, si *common.SuzukiImage) error {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{si.Width, si.Height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < si.Width; x++ {
		for y := 0; y < si.Height; y++ {
			p := si.GetXY(x, y)
			if p == 1 {
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

// SaveContourSliceImage saves a contour (and all child contours) as a PNG.
func SaveContourSliceImage(filename string, c *Contour, width int, height int, flipBook bool, minContourSize int) error {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// naive fill
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.Black)
		}
	}

	colour := 0
	count := 0

	drawContour(img, c, flipBook, minContourSize, colour, &count, filename)
	f, _ := os.Create(filename)
	png.Encode(f, img)
	return nil
}

// drawContour saves a contour to the provided image and then recursively calls to save children to same image
func drawContour(img *image.RGBA, c *Contour, flipBook bool, minContourSize int, colour int, count *int, filename string) error {

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
	if c.BorderType == Outer {
		colour = 0
	}

	// draw contour itself.
	if len(c.Points) > 0 && len(c.Points) > minContourSize {

		colourToUse := colours[colour]
		for _, p := range c.Points {
			img.Set(p.X, p.Y, colourToUse)
		}
		colour++
		if colour >= max {
			colour = 0
		}

		// save new image per contour added...  crazy
		if flipBook {
			fn := fmt.Sprintf("%s-%d.png", filename, *count)
			f, _ := os.Create(fn)
			png.Encode(f, img)
			f.Close()
		}
		*count = *count + 1
	}

	for _, child := range c.Children {
		colour++
		if colour >= max {
			colour = 0
		}
		*count = *count + 1
		drawContour(img, child, flipBook, minContourSize, colour, count, filename)
	}

	return nil
}
