# Borders

This is an implementation of the Suzuki + Abe "Topological Structural Analysis of Digitized Binary Images by Border Following"
( http://pdf.xuebalib.com:1262/xuebalib.com.17233.pdf )

This implementation is based off the original paper as well as being inspired by other papers/implementations/ideas from across the web.

## Description

Borders generates a tree of contours/borders which are generated from a monochrome bitmap. The bitmap needs to be either black/0 for background and white/1 for content to be scanned for borders.
The result is a tree of borders due to borders can contain inner borders (holes), which can in turn contain other borders etc.

The API is NOT stable yet and is in progress.

## Usage

Easiest approach is to have a black and white png file and execute:

```
img := border.LoadImage("test.png")

// contour is root node
contour := border.FindContours(img)   

// this will save contours in contour.png
border.SaveContourSliceImage("contour.png", contour, img.Width, img.Height, false, 0) 
```






