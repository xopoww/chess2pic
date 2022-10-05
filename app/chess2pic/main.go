package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/xopoww/chess2pic/pkg/chess"
	"github.com/xopoww/chess2pic/pkg/pic"
)

var args struct{
	notation string

	input  string
	data   string
	output string
}

func init() {
	if pic.DefaultCollection == nil {
		os.Exit(1)
	}

	flag.StringVar(&args.notation, "notation", "", "notation syntax name")
	flag.StringVar(&args.input, "in", "", "input file name")
	flag.StringVar(&args.data, "data", "", "input text")
	flag.StringVar(&args.output, "out", "chess2pic.png", "output file name")
}

func main() {
	flag.Parse()

	if args.notation == "" {
		fmt.Fprint(os.Stderr, "--notation is required")
		os.Exit(1)
	}

	var r io.RuneReader
	if args.input != "" {
		f, err := os.Open(args.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %q: %s", args.input, err)
			os.Exit(1)
		}
		r = bufio.NewReader(f)
	} else {
		r = strings.NewReader(args.data)
	}

	switch args.notation {
	case "fen":
		pos, err := chess.FEN().Parse(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %s", err)
			os.Exit(1)
		}

		img := image.NewRGBA(pic.DefaultCollection.Board().Bounds())
		pic.DrawPosition(img, pic.DefaultCollection, pos)
		
		f, err := os.Create(args.output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating %q: %s", args.output, err)
			os.Exit(1)
		}
		err = png.Encode(f, img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing png file: %s", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown notation: %q", args.notation)
		os.Exit(1)
	}
}