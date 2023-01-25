package main

import (
	"fmt"

	"github.com/kpfaulkner/borders/border"
)

func displayContour(cont border.Contour) {

	for _, p := range cont.Points {
		fmt.Printf("(%d, %d)\n", p.X, p.Y)
	}

	for _, ch := range cont.Children {
		displayContour(*ch)
	}
}
func main() {
	fmt.Printf("So it begins...\n")

	img, _ := border.LoadImage("test4.png", false)
	cont := border.FindContours(img)
	displayContour(*cont)
}
