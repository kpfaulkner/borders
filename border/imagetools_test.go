package border

import (
	"testing"
)

// TestImageTools tests loading/saving of images.
func TestLoadSaveImage(t *testing.T) {

	testImage, err := LoadImage(`../testimages/unittest2.png`, 1, 1)
	if err != nil {
		t.Errorf("Unable to load test image: %s", err.Error())
	}

	err = SaveImage(`c:\temp\test-again.png`, testImage)
	if err != nil {
		t.Errorf("Unable to save test image: %s", err.Error())
	}
}

func TestSaveContourImage(t *testing.T) {

	testImage, err := LoadImage(`../testimages/unittest2.png`, 1, 1)
	if err != nil {
		t.Errorf("Unable to load test image: %s", err.Error())
	}

	cont, err := FindContours(testImage)
	if err != nil {
		t.Errorf("Unable to find contours: %s", err.Error())
	}
	err = SaveContourSliceImage(`c:\temp\test-conture.png`, cont, testImage.Width, testImage.Height, false, 0)
	if err != nil {
		t.Errorf("Unable to save contour image: %s", err.Error())
	}
}
