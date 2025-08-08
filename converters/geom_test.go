package converters

import (
	"math"
	"testing"

	"github.com/kpfaulkner/borders/border"
)

const (

	// lat/lon degTolerance during conversion calculations
	degTolerance = 0.000001
)

// TestNewSlippyToLatLongConverter tests conversion of slippy co-ords to top left lat/long of box
func TestNewSlippyToLatLongConverter(t *testing.T) {
	testCases := []struct {
		name        string
		slippyX     float64
		slippyY     float64
		scale       int
		expectedLon float64
		expectedLat float64
	}{
		{
			name:        "success",
			slippyX:     1891519.0,
			slippyY:     1285047.0,
			scale:       21,
			expectedLon: 144.700756072,
			expectedLat: -37.569480700,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conv := NewSlippyToLatLongConverter(tc.slippyX, tc.slippyY, tc.scale)
			lon, lat := conv(0, 0)

			if math.Abs(lon-tc.expectedLon) > degTolerance {
				t.Errorf("expected lon %f, got %f", tc.expectedLon, lon)
			}

			if math.Abs(lat-tc.expectedLat) > degTolerance {
				t.Errorf("expected lat %f, got %f", tc.expectedLat, lat)
			}
		})
	}
}

func TestMultiPolygonOnlyConvertContourToPolygon(t *testing.T) {

	testImage, err := border.LoadImage(`../testimages/unittest1.png`, 1, 1)
	if err != nil {
		t.Errorf("Unable to load test image: %s", err.Error())
	}

	cont, err := border.FindContours(testImage)
	if err != nil {
		t.Fatalf("Unable to find contours: %s", err.Error())
	}

	poly, err := ConvertContourToPolygon(cont, 21, true, 0, 0, true)
	if err != nil {
		t.Fatalf("Unable to convert to simple polygon: %s", err.Error())
	}

	if poly.AsText() != "MULTIPOLYGON(((1 1,1 34,33 34,33 28,28 28,27 27,25 27,23 25,23 24,22 23,23 22,23 20,29 14,30 15,30 17,32 17,32 15,30 15,29 14,29 12,31 10,31 7,30 7,29 6,30 5,33 5,33 1,1 1),(23 8,24 7,25 8,25 9,24 10,23 9,23 8),(17 9,18 8,19 9,18 10,17 9)))" {
		t.Errorf("expected polygon to be MULTIPOLYGON(((1 1,1 34,33 34,33 28,28 28,27 27,25 27,23 25,23 24,22 23,23 22,23 20,29 14,30 15,30 17,32 17,32 15,30 15,29 14,29 12,31 10,31 7,30 7,29 6,30 5,33 5,33 1,1 1),(23 8,24 7,25 8,25 9,24 10,23 9,23 8),(17 9,18 8,19 9,18 10,17 9))), got %s", poly.AsText())
	}
}

func TestMultiPolygonConvertContourToPolygon(t *testing.T) {

	testImage, err := border.LoadImage(`../testimages/unittest2.png`, 1, 1)
	if err != nil {
		t.Errorf("Unable to load test image: %s", err.Error())
	}

	cont, err := border.FindContours(testImage)
	if err != nil {
		t.Fatalf("Unable to find contours: %s", err.Error())
	}

	poly, err := ConvertContourToPolygon(cont, 21, true, 0, 0, false)
	if err != nil {
		t.Fatalf("Unable to convert to simple polygon: %s", err.Error())
	}

	if poly.AsText() != "MULTIPOLYGON(((1 1,1 6,4 6,4 4,3 3,3 1,1 1)))" {
		t.Errorf("expected polygon to be MULTIPOLYGON(((1 1,1 6,4 6,4 4,3 3,3 1,1 1))), got %s", poly.AsText())
	}
}

func TestNotSimplifiedMultiPolygonConvertContourToPolygon(t *testing.T) {

	testImage, err := border.LoadImage(`../testimages/unittest2.png`, 1, 1)
	if err != nil {
		t.Errorf("Unable to load test image: %s", err.Error())
	}

	cont, err := border.FindContours(testImage)
	if err != nil {
		t.Fatalf("Unable to find contours: %s", err.Error())
	}

	poly, err := ConvertContourToPolygon(cont, 21, false, 0, 0, false)
	if err != nil {
		t.Fatalf("Unable to convert to simple polygon: %s", err.Error())
	}

	if poly.AsText() != "MULTIPOLYGON(((1 1,1 2,1 3,1 4,1 5,1 6,2 6,3 6,4 6,4 5,4 4,3 3,3 2,3 1,2 1,1 1)))" {
		t.Errorf("expected polygon to be MULTIPOLYGON(((1 1,1 2,1 3,1 4,1 5,1 6,2 6,3 6,4 6,4 5,4 4,3 3,3 2,3 1,2 1,1 1))), got %s", poly.AsText())
	}
}

func TestLatLongToSlippy(t *testing.T) {
	testCases := []struct {
		name            string
		expectedSlippyX float64
		expectedSlippyY float64
		scale           int
		lon             float64
		lat             float64
	}{

		{
			name:            "success",
			expectedSlippyX: 1891519.0,
			expectedSlippyY: 1285047.0,
			scale:           21,
			lon:             144.7007660,
			lat:             -37.5694910,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, y := LatLongToSlippy(tc.lat, tc.lon, tc.scale)
			if x != tc.expectedSlippyX {
				t.Errorf("expected slippyX %f, got %f", tc.expectedSlippyX, x)
			}
			if y != tc.expectedSlippyY {
				t.Errorf("expected slippyY %f, got %f", tc.expectedSlippyY, y)
			}
		})
	}
}
