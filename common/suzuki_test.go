package common

import (
	"image"
	"slices"
	"testing"
)

func createSuzukiImage(t *testing.T, data []int) *SuzukiImage {

	si := NewSuzukiImageFromData(5, 5, data)
	si.SetXY(0, 0, 1)
	v := si.GetXY(0, 0)
	if v != 1 {
		t.Errorf("expected value 1, got %d", v)
	}

	p := image.Point{1, 1}
	si.Set(p, 1)
	if si.Get(p) != 1 {
		t.Errorf("expected value 1, got %d", si.Get(p))
	}

	return si
}

// TestSuzuki tests creation of SuzukiImage.
// Given minimal functionality, just combining the tests for now.
func TestSuzuki(t *testing.T) {

	data := []int{
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0}

	si := createSuzukiImage(t, data)

	getData := si.GetAllData()
	if len(getData) != 25 {
		t.Errorf("expected data length 25, got %d", len(data))
	}

	textSlice := si.DisplayAsText()
	if len(textSlice) != 5 {
		t.Errorf("expected text slice length 5, got %d", len(textSlice))
	}
	if slices.Compare(textSlice, []string{"1 0 0 0 0\n", "0 1 0 0 0\n", "0 0 0 0 0\n", "0 0 0 0 0\n", "0 0 0 0 0\n"}) != 0 {
		t.Errorf("expected text slice to match, got %v", textSlice)
	}

	si2 := createSuzukiImage(t, data)
	if !si.Equals(si2) {
		t.Errorf("expected SuzukiImages to be equal, but they are not")
	}

}
