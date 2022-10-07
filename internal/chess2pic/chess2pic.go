package chess2pic

import (
	"bufio"
	"image"
	"image/gif"
	"image/png"
	"io"
	"os"

	"github.com/andybons/gogif"
	"github.com/xopoww/chess2pic/pkg/chess"
	"github.com/xopoww/chess2pic/pkg/pic"
)

func readerToRuneReader(r io.Reader) io.RuneReader {
	if rs, ok := r.(io.RuneReader); ok {
		return rs
	} else {
		return bufio.NewReader(r)
	}
}

func HandleFEN(in io.Reader, outFile string, col pic.Collection) error {
	rs := readerToRuneReader(in)

	pos, err := chess.FEN().Parse(rs)
	if err != nil {
		return err
	}

	img := pic.DrawPosition(col, pos)
	out, err := os.Create(outFile + ".png")
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, img)
}

func HandlePGN(in io.Reader, outFile string, col pic.Collection) error {
	res, err := chess.ParsePGN(in)
	if err != nil {
		return err
	}

	Debugf("Parsed PGN with %d moves", len(res.Moves))
	Debugf("PGN tags: %#v", res.Tags)

	poss := make([]chess.Position, 0, len(res.Moves) + 1)
	poss = append(poss, res.Start)
	for i, mov := range res.Moves {
		poss = append(poss, chess.Apply(poss[i], mov))
	}

	dst := &gif.GIF{
	}
	quantizer := gogif.MedianCutQuantizer{NumColor: 64}
	for _, pos := range poss {
		img := pic.DrawPosition(col, pos)
		pimg := image.NewPaletted(img.Bounds(), nil)
		quantizer.Quantize(pimg, img.Bounds(), img, image.Point{})
		dst.Image = append(dst.Image, pimg)
		dst.Delay = append(dst.Delay, 100)
	}
	dst.Delay[len(dst.Delay)-1] = 500

	out, err := os.Create(outFile + ".gif")
	if err != nil {
		return err
	}
	defer out.Close()
	return gif.EncodeAll(out, dst)
}