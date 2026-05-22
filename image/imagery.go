package image

import (
	"github.com/kpfaulkner/borders/common"
)

// Erode applies morphological erosion with a square (2*radius+1)×(2*radius+1)
// structuring element. The outermost 1-pixel border is treated as zero, so the
// result's outermost border is always zero. Input is not modified.
//
// Two separable 1-D passes — O(W·H·r) work instead of O(W·H·r²).
// https://en.wikipedia.org/wiki/Erosion_(morphology)
func Erode(img *common.SuzukiImage, radius int) (*common.SuzukiImage, error) {
	width, height := img.Width, img.Height
	tmp := common.NewSuzukiImage(width, height, img.HasPadding())
	out := common.NewSuzukiImage(width, height, img.HasPadding())

	// x is bounded so the window never touches column 0 or width-1, which are
	// treated as permanently zero — that leaves the outer band of the output
	// eroded against the zero edge.
	for y := 0; y < height; y++ {
		for x := radius + 1; x < width-1-radius; x++ {
			v := 1
			for k := x - radius; k <= x+radius; k++ {
				if img.GetXY(k, y) != 1 {
					v = 0
					break
				}
			}
			if v == 1 {
				tmp.SetXY(x, y, 1)
			}
		}
	}

	for y := radius + 1; y < height-1-radius; y++ {
		for x := 0; x < width; x++ {
			v := 1
			for k := y - radius; k <= y+radius; k++ {
				if tmp.GetXY(x, k) != 1 {
					v = 0
					break
				}
			}
			if v == 1 {
				out.SetXY(x, y, 1)
			}
		}
	}

	return out, nil
}

// Dilate applies morphological dilation with a square (2*radius+1)×(2*radius+1)
// structuring element. Out-of-bounds pixels are treated as zero. Input is not
// modified.
//
// Two separable 1-D passes — O(W·H·r) work instead of O(W·H·r²).
// https://en.wikipedia.org/wiki/Dilation_(morphology)
func Dilate(img *common.SuzukiImage, radius int) (*common.SuzukiImage, error) {
	width, height := img.Width, img.Height
	tmp := common.NewSuzukiImage(width, height, img.HasPadding())
	out := common.NewSuzukiImage(width, height, img.HasPadding())

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			lo, hi := x-radius, x+radius
			if lo < 0 {
				lo = 0
			}
			if hi >= width {
				hi = width - 1
			}
			for k := lo; k <= hi; k++ {
				if img.GetXY(k, y) == 1 {
					tmp.SetXY(x, y, 1)
					break
				}
			}
		}
	}

	for y := 0; y < height; y++ {
		lo, hi := y-radius, y+radius
		if lo < 0 {
			lo = 0
		}
		if hi >= height {
			hi = height - 1
		}
		for x := 0; x < width; x++ {
			for k := lo; k <= hi; k++ {
				if tmp.GetXY(x, k) == 1 {
					out.SetXY(x, y, 1)
					break
				}
			}
		}
	}

	return out, nil
}
