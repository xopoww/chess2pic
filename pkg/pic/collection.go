package pic

import (
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"path"

	"github.com/xopoww/chess2pic/pkg/chess"
)

type Image image.Image

// Collection represents a set of Images of all pieces (in both colors) and a board.
// All images must have square dimensions and all piece images must be the same size.
// The size of piece image must not exceed the size of board image divided by 8.
type Collection interface {
	Board() Image
	// Piece returns an image for the specified chess piece.
	// If p.Kind == None, Piece should return nil Image.
	Piece(p chess.Piece) Image
}


type collection struct {
	images [1 + 6 + 6]Image
}

func (col collection) Board() Image {
	return col.images[0]
}

func (col collection) Piece(p chess.Piece) Image {
	if p.Kind == chess.None {
		return nil
	}
	return col.images[1 + int(p.Color * 6) + int(p.Kind - chess.Pawn)]
}


func OpenCollection(dir fs.FS, prefix string) (Collection, error) {
	col := collection{}

	img, err := loadSquareImage(dir, "board.png", prefix)
	if err != nil {
		return col, err
	}
	col.images[0] = img
	bs := img.Bounds().Dx()
	ps := -1

	for color := chess.White; color <= chess.Black; color++ {
		for kind := chess.Pawn; kind <= chess.King; kind++ {
			img, err = loadSquareImage(dir, path.Join(color.Name(), kind.Name() + ".png"), prefix)
			if err != nil {
				return col, err
			}
			col.images[1 + int(color) * 6 + int(kind - chess.Pawn)] = img
			if ps < 0 {
				ps = img.Bounds().Dx()
				if ps * 8 > bs {
					return nil, fmt.Errorf("piece image too big for the board")
				}
			} else {
				if img.Bounds().Dx() != ps {
					return nil, fmt.Errorf("piece images must be the same size")
				}
			}
		}
	}

	return col, nil
}

func loadSquareImage(dir fs.FS, name, prefix string) (Image, error) {
	f, err := dir.Open(path.Join(prefix, name))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", name, err)
	}
	if img.Bounds().Dx() != img.Bounds().Dy() {
		return nil, fmt.Errorf("%s: square image is required", name)
	}
	return img, err
}
