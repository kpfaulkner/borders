package main

import (
	"fmt"
	"image"
	"strings"
)

type SuzukiImage struct {
	Width   int
	Height  int
	data    []int
	dataLen int
}

func NewSuzukiImage(width int, height int) *SuzukiImage {
	si := SuzukiImage{}
	si.Width = width
	si.Height = height
	si.data = make([]int, width*height)
	si.dataLen = width * height // just saves us calculating a lot
	return &si
}

func (si *SuzukiImage) Get(p image.Point) int {
	idx := p.Y*si.Width + p.X
	return si.data[idx]
}

func (si *SuzukiImage) GetXY(x int, y int) int {
	idx := y*si.Width + x

	/*
		// hack...  is somehow possible to get these conditions and blow up.
		// returning 0 is a hack, but results seem good. Need to re-read (for the billionth time)
		// the Suzuki papers etc.
		if idx == -1 || idx < 0 {
			fmt.Printf("BOOM !!! x %d : y %d\n", x, y)
			return 0 // safe?
		}

		// ditto
		if idx >= si.dataLen {
			fmt.Printf("BOOM2 !!! x %d : y %d\n", x, y)
			return 0 // safe?
		}*/

	return si.data[idx]
}

func (si *SuzukiImage) Set(p image.Point, val int) {
	idx := p.Y*si.Width + p.X
	si.data[idx] = val
}

func (si *SuzukiImage) SetXY(x int, y int, val int) {
	idx := y*si.Width + x
	si.data[idx] = val
}

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
