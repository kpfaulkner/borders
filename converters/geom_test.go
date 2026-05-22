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

	if poly.AsText() != "MULTIPOLYGON(((0 0,0 34,34 34,34 28,28 28,27 27,25 27,23 25,23 24,22 23,23 22,23 20,29 14,30 15,30 17,32 17,32 15,30 15,29 14,29 12,31 10,31 7,32 6,34 6,34 0,0 0),(23 8,24 7,25 8,25 9,24 10,23 9,23 8),(17 9,18 8,19 9,18 10,17 9)))" {
		t.Errorf("expected polygon to be MULTIPOLYGON(((0 0,0 34,34 34,34 28,28 28,27 27,25 27,23 25,23 24,22 23,23 22,23 20,29 14,30 15,30 17,32 17,32 15,30 15,29 14,29 12,31 10,31 7,32 6,34 6,34 0,0 0),(23 8,24 7,25 8,25 9,24 10,23 9,23 8),(17 9,18 8,19 9,18 10,17 9))), got %s", poly.AsText())
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

func TestPixelXYToLatLong(t *testing.T) {
	const tol = 1e-9

	testCases := []struct {
		name        string
		pixelX      int64
		pixelY      int64
		scale       int
		expectedLat float64
		expectedLon float64
	}{
		{
			name:        "tile center at scale 0 maps to origin",
			pixelX:      128, // 256/2
			pixelY:      128,
			scale:       0,
			expectedLat: 0,
			expectedLon: 0,
		},
		{
			name:        "globe center at scale 21 maps to origin",
			pixelX:      268435456, // 256 * 2^21 / 2
			pixelY:      268435456,
			scale:       21,
			expectedLat: 0,
			expectedLon: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lat, lon := PixelXYToLatLong(tc.pixelX, tc.pixelY, tc.scale)
			if math.Abs(lat-tc.expectedLat) > tol {
				t.Errorf("expected lat %f, got %f", tc.expectedLat, lat)
			}
			if math.Abs(lon-tc.expectedLon) > tol {
				t.Errorf("expected lon %f, got %f", tc.expectedLon, lon)
			}
		})
	}
}

func TestLatLongToPixelXY(t *testing.T) {
	testCases := []struct {
		name      string
		lat       float64
		lon       float64
		scale     int
		expectedX int64
		expectedY int64
	}{
		{
			name:      "origin at scale 0 is tile center",
			lat:       0,
			lon:       0,
			scale:     0,
			expectedX: 128, // 256/2
			expectedY: 128,
		},
		{
			name:      "origin at scale 21 is globe center",
			lat:       0,
			lon:       0,
			scale:     21,
			expectedX: 268435456, // 256 * 2^21 / 2
			expectedY: 268435456,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, y := LatLongToPixelXY(tc.lat, tc.lon, tc.scale)
			if x != tc.expectedX {
				t.Errorf("expected x %d, got %d", tc.expectedX, x)
			}
			if y != tc.expectedY {
				t.Errorf("expected y %d, got %d", tc.expectedY, y)
			}
		})
	}
}

// TestLatLongToPixelXY_NegativeLongitudeProducesNegativeX is a regression test
// for the uint64→int64 return type change. Previously a longitude < -180 (or
// any value driving x below 0) would silently wrap to a huge positive uint64.
// With int64 the result must remain negative so callers can detect/clamp.
func TestLatLongToPixelXY_NegativeLongitudeProducesNegativeX(t *testing.T) {
	x, _ := LatLongToPixelXY(0, -181, 10)
	if x >= 0 {
		t.Errorf("expected negative x for out-of-range longitude, got %d", x)
	}
}

// TestNewPixelToLatLongConverter verifies the (lat, lon, scale) parameter order
// by passing (0, 0) — the image's top-left pixel — to the returned closure and
// asserting it round-trips to the lat/lon the converter was constructed with.
// A parameter-order swap would produce errors of tens of degrees and trip this
// test loudly.
func TestNewPixelToLatLongConverter(t *testing.T) {
	// Loose enough to absorb pixel-quantisation round-trip error at scale 18
	// (~5e-6 deg/pixel), tight enough that a lat/lon swap (~50+ deg) fails.
	const tol = 1e-4

	testCases := []struct {
		name  string
		lat   float64
		lon   float64
		scale int
	}{
		{name: "melbourne", lat: -37.5694910, lon: 144.7007660, scale: 21},
		{name: "london", lat: 51.5074, lon: -0.1278, scale: 18},
		{name: "southern hemisphere, negative longitude", lat: -10.0, lon: -75.0, scale: 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conv := NewPixelToLatLongConverter(tc.lat, tc.lon, tc.scale)

			// (0, 0) addresses the image's top-left pixel — the converter should
			// return (lon, lat) matching what we constructed it with.
			lon, lat := conv(0, 0)
			if math.Abs(lon-tc.lon) > tol {
				t.Errorf("top-left lon: expected %f, got %f", tc.lon, lon)
			}
			if math.Abs(lat-tc.lat) > tol {
				t.Errorf("top-left lat: expected %f, got %f", tc.lat, lat)
			}
		})
	}
}

// TestLatLongPixelRoundTrip asserts the two functions are near-inverses within
// pixel-quantisation error.
func TestLatLongPixelRoundTrip(t *testing.T) {
	const tol = 1e-3 // degrees; pixel quantisation is much coarser than this at low zooms

	testCases := []struct {
		name  string
		lat   float64
		lon   float64
		scale int
	}{
		{name: "melbourne", lat: -37.5694910, lon: 144.7007660, scale: 21},
		{name: "london", lat: 51.5074, lon: -0.1278, scale: 18},
		{name: "southern hemisphere, negative longitude", lat: -10.0, lon: -75.0, scale: 15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, y := LatLongToPixelXY(tc.lat, tc.lon, tc.scale)
			lat, lon := PixelXYToLatLong(x, y, tc.scale)
			if math.Abs(lat-tc.lat) > tol {
				t.Errorf("lat round-trip: expected %f, got %f", tc.lat, lat)
			}
			if math.Abs(lon-tc.lon) > tol {
				t.Errorf("lon round-trip: expected %f, got %f", tc.lon, lon)
			}
		})
	}
}
