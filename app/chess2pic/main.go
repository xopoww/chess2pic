package main

import (
	"image"
	"image/png"
	"os"

	"github.com/xopoww/chess2pic/pkg/chess"
	"github.com/xopoww/chess2pic/pkg/pic"
)

func main() {
	if pic.DefaultCollection == nil {
		os.Exit(1)
	}

	img := image.NewRGBA(pic.DefaultCollection.Board().Bounds())
	pic.DrawPosition(img, pic.DefaultCollection, chess.StartingPosition())
	
	f, err := os.Create("chess.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}