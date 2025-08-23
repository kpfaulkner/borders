// Package converters is used to convert/translate the generated border into different format such
// GIS coordinates (latitude/longitude), GeoJSON, Slippy coords.
//
// The key functions to use are:
//   NewSlippyToLatLongConverter : This returns a function that will generate Slippy coordinates into
//   GIS coordinates. This is useful if the input image is a bitmap of Slippy Coords ( see https://en.wikipedia.org/wiki/Tiled_web_map )
//   but the final result is required to be a GeoJSON geometry which can be used in GIS applications.
//   See examples/slippy-lat-log
//
//   NewPixelToLatLongConverter is similarly used if the input image is a map and the output is a GeoJSON geometry.
//   This is similar to NewSlippyToLatLongConverter but more "fine grain".
//
//   ConvertContourToPolygon is a more generic function that takes a generated Contour and converts to a
//   Geometry. This will be used in combination with NewSlippyToLatLongConverter or NewPixelToLatLongConverter

package converters
