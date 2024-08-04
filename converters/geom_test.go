package converters

import (
	"math"
	"testing"
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
