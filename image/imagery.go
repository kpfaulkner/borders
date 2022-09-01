package image

import "github.com/kpfaulkner/borders/border"

// Erode the suzuki image, based on Morphological Erosion
// https://en.wikipedia.org/wiki/Erosion_(morphology)
// THIS IS GOOD... AT LEAST SAME AS BILT
func Erode(img *border.SuzukiImage, radius int) (*border.SuzukiImage, error) {

	img2 := border.NewSuzukiImage(img.Width, img.Height)
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {

			// if x == 0 or y == 0 or x == img.Width-1 or y == img.Height-1 then its an edge, and set it to 0.
			if x == 0 || y == 0 || x == img.Width-1 || y == img.Height-1 {
				img.SetXY(x, y, 0)
				continue
			}

			// for each pixel, check if all pixels within radius are 1
			// if not, set to 0
			if img.GetXY(x, y) == 1 {
				// check if all pixels within radius are 1
				// if not, set to 0
				if !checkErodeRadius(img, x, y, img.Width, img.Height, radius) {
					img2.SetXY(x, y, 0)
				} else {
					img2.SetXY(x, y, 1)
				}
			}
		}
	}
	return img2, nil
}

func checkErodeRadius(img *border.SuzukiImage, x int, y int, width int, height int, radius int) bool {
	for i := -radius; i <= radius; i++ {
		for j := -radius; j <= radius; j++ {
			if x+i < 0 || y+j < 0 || x+i >= width || y+j >= height {
				continue // out of bounds.
			}
			if img.GetXY(x+i, y+j) != 1 {
				return false
			}
		}
	}
	return true
}

// Dilate the suzuki image, based on Morphological Dilation
// https://en.wikipedia.org/wiki/Dilation_(morphology)
// WORKS
func Dilate(img *border.SuzukiImage, radius int) (*border.SuzukiImage, error) {
	img2 := border.NewSuzukiImage(img.Width, img.Height)

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {

			// for each pixel, check if any pixels within radius are 1
			// if not, set to 0
			if img.GetXY(x, y) == 1 {
				dilateRadiusAroundPoint(img, img2, x, y, img.Width, img.Height, radius)
			}
		}
	}
	return img2, nil
}

func dilateRadiusAroundPoint(img *border.SuzukiImage, img2 *border.SuzukiImage, x int, y int, width int, height int, radius int) {
	for i := -radius; i <= radius; i++ {
		for j := -radius; j <= radius; j++ {
			if x+i < 0 || y+j < 0 || x+i >= width || y+j >= height {
				continue // out of bounds.
			}
			img2.SetXY(x+i, y+j, 1)
			/*
				if img.GetXY(x+i, y+j) == 1 {
					return true
				} */
		}
	}
}
