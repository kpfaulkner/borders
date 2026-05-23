package border

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"

	"github.com/kpfaulkner/borders/common"
	image2 "github.com/kpfaulkner/borders/image"
)

// LoadImage loads a PNG and returns a SuzukiImage. Currently restricted to PNG but will eventually expand
// to include other formats.
//
// Erode parameter forces the eroding of the image before converting to a SuzukiImage.
// See https://en.wikipedia.org/wiki/Erosion_(morphology) for explanation
//
// Dilate parameter forces the dilating of the image before converting to a SuzukiImage. Likewise, see
// https://en.wikipedia.org/wiki/Dilation_(morphology) for explanation.
//
// The combination of Erode and Dilate helps remove any "spikes" that may appear in the generated boundary.
// erode and dilate will usually be 0 (none) or 1 (single pixel spikes)
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

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	isForeground := newForegroundTester(img)

	// If any pixels on the edges are populated, then we need to pad this out by 1 pixel on each side.
	// This will be reversed later.
	requirePadding := edgeHasForeground(width, height, isForeground)

	// need border to be black. Pad edges with 1 black pixel
	si := common.NewSuzukiImage(width, height, requirePadding)

	paddingOffset := 0
	if requirePadding {
		paddingOffset = 1
	}

	// Walk in image-local 0-based coords; isForeground handles concrete-type Pix access.
	// Background pixels are skipped — the SuzukiImage backing slice is already zero.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isForeground(x, y) {
				si.SetXY(x+paddingOffset, y+paddingOffset, 1)
			}
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

// newForegroundTester returns a closure that reports whether the pixel at the
// given image-local 0-based coordinate is foreground (non-black).
//
// The slow generic path uses image.Image.At(x, y).RGBA() per pixel — interface
// dispatch on every call plus, for some image types, a per-call color allocation.
// Direct .Pix access for the concrete types image/png actually returns is
// roughly an order of magnitude faster on large images.
//
// Semantics match the original At/RGBA path: a pixel is foreground iff its
// alpha-premultiplied RGB is non-zero. Fully transparent pixels are background
// regardless of their underlying RGB.
func newForegroundTester(img image.Image) func(x, y int) bool {
	switch im := img.(type) {
	case *image.NRGBA:
		pix, stride := im.Pix, im.Stride
		return func(x, y int) bool {
			off := y*stride + x*4
			// Non-premultiplied: also require non-zero alpha so transparent pixels are background.
			return pix[off+3] != 0 && (pix[off] != 0 || pix[off+1] != 0 || pix[off+2] != 0)
		}
	case *image.RGBA:
		pix, stride := im.Pix, im.Stride
		return func(x, y int) bool {
			off := y*stride + x*4
			// Premultiplied — non-zero RGB already implies non-zero alpha.
			return pix[off] != 0 || pix[off+1] != 0 || pix[off+2] != 0
		}
	case *image.Gray:
		pix, stride := im.Pix, im.Stride
		return func(x, y int) bool {
			return pix[y*stride+x] != 0
		}
	case *image.Paletted:
		// Resolve each palette entry once — the inner loop becomes a single byte lookup.
		isFG := make([]bool, len(im.Palette))
		for i, c := range im.Palette {
			r, g, b, _ := c.RGBA()
			isFG[i] = r != 0 || g != 0 || b != 0
		}
		pix, stride := im.Pix, im.Stride
		return func(x, y int) bool {
			return isFG[pix[y*stride+x]]
		}
	default:
		// Fallback for any uncommon image type — matches the original semantics.
		minX, minY := img.Bounds().Min.X, img.Bounds().Min.Y
		return func(x, y int) bool {
			r, g, b, _ := img.At(x+minX, y+minY).RGBA()
			return r != 0 || g != 0 || b != 0
		}
	}
}

// edgeHasForeground reports whether any pixel along the image's outer edge is
// foreground. Used to decide whether to pad the SuzukiImage with a 1-pixel zero border.
func edgeHasForeground(width, height int, isForeground func(x, y int) bool) bool {
	// left and right columns
	for y := 0; y < height; y++ {
		if isForeground(0, y) || isForeground(width-1, y) {
			return true
		}
	}
	// top and bottom rows
	for x := 0; x < width; x++ {
		if isForeground(x, 0) || isForeground(x, height-1) {
			return true
		}
	}
	return false
}

// SaveImage saves a SuzukiImage as a PNG to given filename.
// Although currently only PNG, will make this more generic in the future.
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

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// SaveContourSliceImage saves a contour (and all child contours) as a PNG.
// Width and height are the dimensions of the image to save.
//
// flipBook is a bool to indicate that when each child contour is added to the image, the image should be
// saved to a new file with the filename suffixed with the count of contours added so far. This is useful
// for debugging and visualising the contours as they are added.
//
// minContourSize indicates if minimum number of points that make up a contour. If contour contains fewer, then
// do NOT save.
func SaveContourSliceImage(filename string, c *Contour, width int, height int, flipBook bool, minContourSize int) (err error) {
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

	if err := drawContour(img, c, flipBook, minContourSize, colour, &count, filename); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	return png.Encode(f, img)
}

// drawContour saves a contour to the provided image and then recursively calls to save children to same image.
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
			f, err := os.Create(fn)
			if err != nil {
				return err
			}
			encErr := png.Encode(f, img)
			if cErr := f.Close(); encErr == nil {
				encErr = cErr
			}
			if encErr != nil {
				return encErr
			}
		}
		*count = *count + 1
	}

	for _, child := range c.Children {
		colour++
		if colour >= max {
			colour = 0
		}
		*count = *count + 1
		if err := drawContour(img, child, flipBook, minContourSize, colour, count, filename); err != nil {
			return err
		}
	}

	return nil
}
