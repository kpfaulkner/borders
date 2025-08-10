package common

import (
	"fmt"
	"image"
	"strings"
)

// SuzukiImage is the basic structure we use to define an image when trying to find contours.
type SuzukiImage struct {
	Width   int
	Height  int
	data    []int
	dataLen int
}

// NewSuzukiImage creates a new SuzukiImage of specific dimensions.
func NewSuzukiImage(width int, height int) *SuzukiImage {
	si := SuzukiImage{}
	si.Width = width
	si.Height = height
	si.data = make([]int, width*height)
	si.dataLen = width * height // just saves us calculating a lot
	return &si
}

func NewSuzukiImageFromData(width int, height int, data []int) *SuzukiImage {
	si := NewSuzukiImage(width, height)
	si.data = data[:]
	si.dataLen = len(data)
	return si
}

// Get returns the value of a given point
func (si *SuzukiImage) GetAllData() []int {
	return si.data
}

// Get returns the value of a given point
func (si *SuzukiImage) Get(p image.Point) int {
	idx := p.Y*si.Width + p.X
	return si.data[idx]
}

// GetXY returns the value of a given x/y
func (si *SuzukiImage) GetXY(x int, y int) int {
	idx := y*si.Width + x
	return si.data[idx]
}

// Set sets the value at a given point
func (si *SuzukiImage) Set(p image.Point, val int) {
	idx := p.Y*si.Width + p.X
	si.data[idx] = val
}

// SetXY sets the value at a given x/y
func (si *SuzukiImage) SetXY(x int, y int, val int) {
	idx := y*si.Width + x
	si.data[idx] = val
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
