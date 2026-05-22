package common

import (
	"fmt"
	"image"
	"strings"
)

// SuzukiImage is the basic structure we use to define an image when trying to find contours.
//
// Storage is int32 per cell: each cell holds either 0/1 (background/foreground)
// or a signed contour id assigned during border following. Contour ids are
// bounded by the pixel count, so int32 is sufficient for any image up to ~2^31
// pixels and halves memory vs. int on 64-bit platforms.
type SuzukiImage struct {
	Width   int
	Height  int
	data    []int32
	dataLen int

	// Indicates if a 1 pixel padding has been applied to around the image.
	// This helps with imagery where it goes RIGHT up to the edge.
	hasPadding bool
}

// NewSuzukiImage creates a new SuzukiImage of specific dimensions.
func NewSuzukiImage(width int, height int, hasPadding bool) *SuzukiImage {
	si := SuzukiImage{}
	padding := 0
	if hasPadding {
		padding = 2
	}
	si.Width = width + padding
	si.Height = height + padding
	si.data = make([]int32, si.Width*si.Height)
	si.dataLen = si.Width * si.Height // just saves us calculating a lot
	si.hasPadding = hasPadding
	return &si
}

func NewSuzukiImageFromData(width int, height int, hasPadding bool, data []int) *SuzukiImage {
	si := NewSuzukiImage(width, height, hasPadding)
	if hasPadding {
		// Copy each row of the unpadded source into the interior of the padded
		// buffer, leaving the 1-pixel zero border intact.
		for y := 0; y < height; y++ {
			dst := (y+1)*si.Width + 1
			src := y * width
			for i := 0; i < width; i++ {
				si.data[dst+i] = int32(data[src+i])
			}
		}
	} else {
		for i, v := range data {
			si.data[i] = int32(v)
		}
	}
	return si
}

// GetAllData returns a copy of the underlying data as []int.
func (si *SuzukiImage) GetAllData() []int {
	out := make([]int, len(si.data))
	for i, v := range si.data {
		out[i] = int(v)
	}
	return out
}

// Get returns the value of a given point
func (si *SuzukiImage) Get(p image.Point) int {
	return int(si.data[p.Y*si.Width+p.X])
}

// GetXY returns the value of a given x/y
func (si *SuzukiImage) GetXY(x int, y int) int {
	return int(si.data[y*si.Width+x])
}

// Set sets the value at a given point
func (si *SuzukiImage) Set(p image.Point, val int) {
	si.data[p.Y*si.Width+p.X] = int32(val)
}

// SetXY sets the value at a given x/y
func (si *SuzukiImage) SetXY(x int, y int, val int) {
	si.data[y*si.Width+x] = int32(val)
}

func (si *SuzukiImage) HasPadding() bool {
	return si.hasPadding
}

// DisplayAsText generates a string of a given image. This is purely used for debugging SMALL images
func (si *SuzukiImage) DisplayAsText() []string {
	s := []string{}
	for y := 0; y < si.Height; y++ {
		ss := si.data[y*si.Width : (y*si.Width + si.Width)]
		t := []string{}
		for _, i := range ss {
			t = append(t, fmt.Sprintf("%d", i))
		}
		s = append(s, strings.Join(t, " ")+"\n")
	}

	return s
}

// Equals checks if two SuzukiImages are equal.
func (si *SuzukiImage) Equals(other *SuzukiImage) bool {
	if si.Width != other.Width || si.Height != other.Height {
		return false
	}

	for i := 0; i < si.dataLen; i++ {
		if si.data[i] != other.data[i] {
			return false
		}
	}
	return true
}
