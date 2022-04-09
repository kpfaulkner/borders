package main

import (
	"fmt"
	"time"

	"github.com/kpfaulkner/borders/border"
)

func main() {
	fmt.Printf("So it begins...\n")

	img := border.LoadImage("image1.png")
	//img := border.LoadImage("tiny.png")
	//img := border.LoadImage("big-test-image.png")

	start := time.Now()
	cont := border.FindContours(img)
	fmt.Printf("finding took %d ms\n", time.Now().Sub(start).Milliseconds())

	//saveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0, false)
	//border.SaveContourSliceImage("contour.png", cont, img.Width, img.Height, false, 0)
	border.SaveContourSliceImage("c:/temp/contour/contour", cont, img.Width, img.Height, true, 0)

	/*
		for _, c := range contours {
			fmt.Printf("%d %d : %d : %+v : %d\n", c.Id, c.ParentId, c.BorderType, c.ConflictingContours, len(c.Children))
		}

		fmt.Printf("Num contours are %d\n", len(cont))

	*/
}
