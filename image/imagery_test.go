package image

import (
	"image"
	"testing"

	"github.com/kpfaulkner/borders/common"
)

// TestErode tests eroding of Suzuki Image
func TestErode(t *testing.T) {
	testCases := []struct {
		name   string
		width  int
		height int
		radius int

		// holes (ie pixels set to 0) in original image.
		holes []image.Point

		// holes (ie pixels set to 0) in eroded image.
		expectedResultData []int

		expectErr bool
	}{
		{
			name:               "success with solid image",
			width:              10,
			height:             10,
			radius:             1,
			holes:              []image.Point{},
			expectedResultData: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:               "success with single pixel missing",
			width:              10,
			height:             10,
			radius:             1,
			holes:              []image.Point{{5, 1}},
			expectedResultData: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},

		{
			name:               "success with two pixels missing",
			width:              10,
			height:             10,
			radius:             1,
			holes:              []image.Point{{5, 1}, {8, 8}},
			expectedResultData: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			incomingImage := createSuzukiImage(tc.width, tc.height, tc.holes)
			resultImage, err := Erode(incomingImage, tc.radius)

			// expected error but didn't get it
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			// didn't expect error but got one
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// expected and received error
			if tc.expectErr && err != nil {
				return
			}

			expectedImage := common.NewSuzukiImageFromData(tc.width, tc.height, tc.expectedResultData)
			if !resultImage.Equals(expectedImage) {
				t.Errorf("result image differs from expected")
			}
		})
	}
}

// TestDilate tests dilating of Suzuki Image
func TestDilate(t *testing.T) {
	testCases := []struct {
		name   string
		width  int
		height int
		radius int

		imageData []int

		// holes (ie pixels set to 0) in eroded image.
		expectedResultData []int

		expectErr bool
	}{
		{
			name:               "success with internal hole",
			width:              11,
			height:             11,
			radius:             1,
			imageData:          []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedResultData: []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			incomingImage := common.NewSuzukiImageFromData(tc.width, tc.height, tc.imageData)
			resultImage, err := Dilate(incomingImage, tc.radius)

			// expected error but didn't get it
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			// didn't expect error but got one
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// expected and received error
			if tc.expectErr && err != nil {
				return
			}

			expectedImage := common.NewSuzukiImageFromData(tc.width, tc.height, tc.expectedResultData)
			if !resultImage.Equals(expectedImage) {
				t.Errorf("result image differs from expected")
			}
		})
	}
}

// createSuzukiImage creates a dummy SuzukiImage for test purposes.
func createSuzukiImage(width int, height int, holes []image.Point) *common.SuzukiImage {
	si := createFullSuzukiImage(width, height)
	for _, p := range holes {
		si.SetXY(p.X, p.Y, 0)
	}
	return si
}

// createFullSuzukiImage creates a SuzukiImage where all pixels are populated (with 1).
func createFullSuzukiImage(width int, height int) *common.SuzukiImage {
	si := common.NewSuzukiImage(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			si.SetXY(x, y, 1)
		}
	}
	return si
}
