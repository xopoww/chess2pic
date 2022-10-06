package pic

import (
	"image"
	"image/draw"

	"github.com/xopoww/chess2pic/pkg/chess"
)

func DrawPosition(dst draw.Image, col Collection, pos chess.Position) {
	if dst.Bounds() != col.Board().Bounds() {
		panic("invalid dst bounds")
	}
	draw.Draw(dst, dst.Bounds(), col.Board(), image.Pt(0, 0), draw.Over)

	bs := dst.Bounds().Dx()
	ss := bs / 8
	ps := col.Piece(chess.Piece{Color: chess.White, Kind: chess.Pawn}).Bounds().Dx()
	off := (ss - ps) / 2

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			sq := chess.MustNewSquare(file, rank)
			p := pos.Get(sq)
			if p.Kind == chess.None {
				continue
			}

			img := col.Piece(p)
			x := file * ss + off
			y := (7 - rank) * ss + off
			draw.Draw(dst, image.Rect(x, y, x + ps, y + ps), img, image.Pt(0, 0), draw.Over)
		}
	}
}