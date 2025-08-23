// Package border implements the core border/boundary detection of an image.
// This is an implementation of the Suzuki + Abe "Topological Structural Analysis of Digitized Binary Images
// by Border Following" ( http://pdf.xuebalib.com:1262/xuebalib.com.17233.pdf ).
//
// This is based off original work based off the paper as well as inspired by other papers and implementations.
//
// The primary function supplied in this package is the FindContours function. This takes a SuzukiImage and
// returns a Contour instance.
//
// A Contour contains all the points of a border. Note: the border is not just the outer border but can
// contain "holes" and sub-borders.
//
// This package also converts a PNG image (will extend to other formats in the future) and generates
// a SuzukiImage instance. A SuzukiImage effectively a wrapper to a byte array with helper functions used
// to make border dection easier.
package border
