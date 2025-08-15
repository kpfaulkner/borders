package border

import (
	"fmt"
	"image"
	"slices"
	"testing"

	"github.com/kpfaulkner/borders/common"
)

// TestFindContour tests conversion of slippy co-ords to top left lat/long of box
func TestFindContour(t *testing.T) {

	testImage, err := LoadImage(`../testimages/unittest1.png`, 1, 1)
	if err != nil {
		t.Fatalf("Unable to load test image: %s", err.Error())
	}
	testCases := []struct {
		name   string
		width  int
		height int
		radius int

		imageData       []int
		preCreatedImage *common.SuzukiImage

		// holes (ie pixels set to 0) in eroded image.
		expectedContourPoints []image.Point

		hasPadding bool
		expectErr  bool
	}{
		{
			name:                  "success with internal hole",
			width:                 10,
			height:                10,
			radius:                1,
			imageData:             []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			hasPadding:            false,
			expectedContourPoints: []image.Point{image.Point{X: 1, Y: 1}, image.Point{X: 1, Y: 2}, image.Point{X: 1, Y: 3}, image.Point{X: 1, Y: 4}, image.Point{X: 1, Y: 5}, image.Point{X: 1, Y: 6}, image.Point{X: 2, Y: 6}, image.Point{X: 3, Y: 6}, image.Point{X: 4, Y: 6}, image.Point{X: 4, Y: 5}, image.Point{X: 4, Y: 4}, image.Point{X: 3, Y: 3}, image.Point{X: 3, Y: 2}, image.Point{X: 3, Y: 1}, image.Point{X: 2, Y: 1}},
		},
		{
			name:                  "success with colliding with another contour",
			width:                 testImage.Width,
			height:                testImage.Height,
			radius:                1,
			hasPadding:            testImage.HasPadding(),
			preCreatedImage:       testImage,
			expectedContourPoints: []image.Point{image.Point{X: 0, Y: 0}, image.Point{X: 0, Y: 1}, image.Point{X: 0, Y: 2}, image.Point{X: 0, Y: 3}, image.Point{X: 0, Y: 4}, image.Point{X: 0, Y: 5}, image.Point{X: 0, Y: 6}, image.Point{X: 0, Y: 7}, image.Point{X: 0, Y: 8}, image.Point{X: 0, Y: 9}, image.Point{X: 0, Y: 10}, image.Point{X: 0, Y: 11}, image.Point{X: 0, Y: 12}, image.Point{X: 0, Y: 13}, image.Point{X: 0, Y: 14}, image.Point{X: 0, Y: 15}, image.Point{X: 0, Y: 16}, image.Point{X: 0, Y: 17}, image.Point{X: 0, Y: 18}, image.Point{X: 0, Y: 19}, image.Point{X: 0, Y: 20}, image.Point{X: 0, Y: 21}, image.Point{X: 0, Y: 22}, image.Point{X: 0, Y: 23}, image.Point{X: 0, Y: 24}, image.Point{X: 0, Y: 25}, image.Point{X: 0, Y: 26}, image.Point{X: 0, Y: 27}, image.Point{X: 0, Y: 28}, image.Point{X: 0, Y: 29}, image.Point{X: 0, Y: 30}, image.Point{X: 0, Y: 31}, image.Point{X: 0, Y: 32}, image.Point{X: 0, Y: 33}, image.Point{X: 0, Y: 34}, image.Point{X: 1, Y: 34}, image.Point{X: 2, Y: 34}, image.Point{X: 3, Y: 34}, image.Point{X: 4, Y: 34}, image.Point{X: 5, Y: 34}, image.Point{X: 6, Y: 34}, image.Point{X: 7, Y: 34}, image.Point{X: 8, Y: 34}, image.Point{X: 9, Y: 34}, image.Point{X: 10, Y: 34}, image.Point{X: 11, Y: 34}, image.Point{X: 12, Y: 34}, image.Point{X: 13, Y: 34}, image.Point{X: 14, Y: 34}, image.Point{X: 15, Y: 34}, image.Point{X: 16, Y: 34}, image.Point{X: 17, Y: 34}, image.Point{X: 18, Y: 34}, image.Point{X: 19, Y: 34}, image.Point{X: 20, Y: 34}, image.Point{X: 21, Y: 34}, image.Point{X: 22, Y: 34}, image.Point{X: 23, Y: 34}, image.Point{X: 24, Y: 34}, image.Point{X: 25, Y: 34}, image.Point{X: 26, Y: 34}, image.Point{X: 27, Y: 34}, image.Point{X: 28, Y: 34}, image.Point{X: 29, Y: 34}, image.Point{X: 30, Y: 34}, image.Point{X: 31, Y: 34}, image.Point{X: 32, Y: 34}, image.Point{X: 33, Y: 34}, image.Point{X: 34, Y: 34}, image.Point{X: 34, Y: 33}, image.Point{X: 34, Y: 32}, image.Point{X: 34, Y: 31}, image.Point{X: 34, Y: 30}, image.Point{X: 34, Y: 29}, image.Point{X: 34, Y: 28}, image.Point{X: 33, Y: 28}, image.Point{X: 32, Y: 28}, image.Point{X: 31, Y: 28}, image.Point{X: 30, Y: 28}, image.Point{X: 29, Y: 28}, image.Point{X: 28, Y: 28}, image.Point{X: 27, Y: 27}, image.Point{X: 26, Y: 27}, image.Point{X: 25, Y: 27}, image.Point{X: 24, Y: 26}, image.Point{X: 23, Y: 25}, image.Point{X: 23, Y: 24}, image.Point{X: 22, Y: 23}, image.Point{X: 23, Y: 22}, image.Point{X: 23, Y: 21}, image.Point{X: 23, Y: 20}, image.Point{X: 24, Y: 19}, image.Point{X: 25, Y: 18}, image.Point{X: 26, Y: 17}, image.Point{X: 27, Y: 16}, image.Point{X: 28, Y: 15}, image.Point{X: 29, Y: 14}, image.Point{X: 30, Y: 15}, image.Point{X: 30, Y: 16}, image.Point{X: 30, Y: 17}, image.Point{X: 31, Y: 17}, image.Point{X: 32, Y: 17}, image.Point{X: 32, Y: 16}, image.Point{X: 32, Y: 15}, image.Point{X: 31, Y: 15}, image.Point{X: 30, Y: 15}, image.Point{X: 29, Y: 14}, image.Point{X: 29, Y: 13}, image.Point{X: 29, Y: 12}, image.Point{X: 30, Y: 11}, image.Point{X: 31, Y: 10}, image.Point{X: 31, Y: 9}, image.Point{X: 31, Y: 8}, image.Point{X: 31, Y: 7}, image.Point{X: 32, Y: 6}, image.Point{X: 33, Y: 6}, image.Point{X: 34, Y: 6}, image.Point{X: 34, Y: 5}, image.Point{X: 34, Y: 4}, image.Point{X: 34, Y: 3}, image.Point{X: 34, Y: 2}, image.Point{X: 34, Y: 1}, image.Point{X: 34, Y: 0}, image.Point{X: 33, Y: 0}, image.Point{X: 32, Y: 0}, image.Point{X: 31, Y: 0}, image.Point{X: 30, Y: 0}, image.Point{X: 29, Y: 0}, image.Point{X: 28, Y: 0}, image.Point{X: 27, Y: 0}, image.Point{X: 26, Y: 0}, image.Point{X: 25, Y: 0}, image.Point{X: 24, Y: 0}, image.Point{X: 23, Y: 0}, image.Point{X: 22, Y: 0}, image.Point{X: 21, Y: 0}, image.Point{X: 20, Y: 0}, image.Point{X: 19, Y: 0}, image.Point{X: 18, Y: 0}, image.Point{X: 17, Y: 0}, image.Point{X: 16, Y: 0}, image.Point{X: 15, Y: 0}, image.Point{X: 14, Y: 0}, image.Point{X: 13, Y: 0}, image.Point{X: 12, Y: 0}, image.Point{X: 11, Y: 0}, image.Point{X: 10, Y: 0}, image.Point{X: 9, Y: 0}, image.Point{X: 8, Y: 0}, image.Point{X: 7, Y: 0}, image.Point{X: 6, Y: 0}, image.Point{X: 5, Y: 0}, image.Point{X: 4, Y: 0}, image.Point{X: 3, Y: 0}, image.Point{X: 2, Y: 0}, image.Point{X: 1, Y: 0}, image.Point{X: 29, Y: 6}, image.Point{X: 30, Y: 5}, image.Point{X: 31, Y: 5}, image.Point{X: 32, Y: 6}, image.Point{X: 31, Y: 7}, image.Point{X: 30, Y: 7}, image.Point{X: 23, Y: 8}, image.Point{X: 24, Y: 7}, image.Point{X: 25, Y: 8}, image.Point{X: 25, Y: 9}, image.Point{X: 24, Y: 10}, image.Point{X: 23, Y: 9}, image.Point{X: 17, Y: 9}, image.Point{X: 18, Y: 8}, image.Point{X: 19, Y: 9}, image.Point{X: 18, Y: 10}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var incomingImage *common.SuzukiImage
			if tc.preCreatedImage != nil {
				incomingImage = tc.preCreatedImage
			} else {
				incomingImage = common.NewSuzukiImageFromData(tc.width, tc.height, tc.hasPadding, tc.imageData)
			}
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

			comparePoints := func(a, b image.Point) int {
				if a.X < b.X {
					return -1 // a is less than b
				}
				if a.X > b.X {
					return 1 // a is greater than b
				}

				if a.Y < b.Y {
					return -1 // a is less than b
				}

				if a.Y > b.Y {
					return 1 // a is greater than b
				}

				return 0
			}

			if slices.CompareFunc(cont.GetAllPoints(), tc.expectedContourPoints, comparePoints) != 0 {
				fmt.Printf("got %#v\n", cont.GetAllPoints())
				t.Errorf("expected contour points %v, got %v", tc.expectedContourPoints, cont.GetAllPoints())
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
	si := common.NewSuzukiImage(width, height, false)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			si.SetXY(x, y, 1)
		}
	}
	return si
}
