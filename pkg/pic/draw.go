package pic

import (
	// 	"image"

	"image"
	"image/draw"

	"github.com/xopoww/chess2pic/pkg/chess"
)

func DrawPosition(dst draw.Image, col Collection, pos chess.Position) {
	draw.Draw(dst, dst.Bounds(), col.Board(), image.Pt(0, 0), draw.Over)
}