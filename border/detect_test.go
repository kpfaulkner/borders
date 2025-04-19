package border

import (
	"fmt"
	"image"
	"testing"

	"github.com/kpfaulkner/borders/common"
)

// TestFindContour tests conversion of slippy co-ords to top left lat/long of box
func TestFindContour(t *testing.T) {
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
			width:              10,
			height:             10,
			radius:             1,
			imageData:          []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expectedResultData: []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			incomingImage := common.NewSuzukiImageFromData(tc.width, tc.height, tc.imageData)
			cont, err := FindContours(incomingImage)
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// expected and received error
			if tc.expectErr && err != nil {
				return
			}
			fmt.Printf("contour %+v\n", cont)

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
