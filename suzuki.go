package main

import "image"

type SuzukiImage struct {
	Width  int
	Height int
	data   []byte
}

func NewSuzukiImage(width int, height int) *SuzukiImage {
	si := SuzukiImage{}
	si.Width = width
	si.Height = height
	si.data = make([]byte, width*height)
	return &si
}

func (si *SuzukiImage) Get(p image.Point) byte {
	idx := p.Y*si.Width + p.X
	return si.data[idx]
}

func (si *SuzukiImage) Set(p image.Point, val byte) {
	idx := p.Y*si.Width + p.X
	si.data[idx] = val
}
