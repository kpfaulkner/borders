package main

import "image"

type SuzukiImage struct {
	Width  int
	Height int
	data   []int
}

func NewSuzukiImage(width int, height int) *SuzukiImage {
	si := SuzukiImage{}
	si.Width = width
	si.Height = height
	si.data = make([]int, width*height)
	return &si
}

func (si *SuzukiImage) Get(p image.Point) int {
	idx := p.Y*si.Width + p.X
	return si.data[idx]
}

func (si *SuzukiImage) GetXY(x int, y int) int {
	idx := y*si.Width + x
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
