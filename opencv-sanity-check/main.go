package main

import (
	"fmt"
	"time"

	"github.com/kpfaulkner/borders/border"
)

func displayContour(cont border.Contour) {

	for _, p := range cont.Points {
		fmt.Printf("(%d,%d)\n", p.X, p.Y)
	}

	for _, ch := range cont.Children {
		displayContour(*ch)
	}
}
func main() {
	img, _ := border.LoadImage("test.png", 0, 0)

	start := time.Now()
	cont, _ := border.FindContours(img)
	fmt.Printf("took %d ms\n", time.Since(start).Milliseconds())
	displayContour(*cont)
}
