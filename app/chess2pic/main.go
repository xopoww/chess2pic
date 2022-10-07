package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xopoww/chess2pic/internal/chess2pic"
	"github.com/xopoww/chess2pic/pkg/chess"
	"github.com/xopoww/chess2pic/pkg/pic"
)

var args struct {
	notation string

	input  string
	data   string
	output string

	from string
}

func init() {
	if pic.DefaultCollection == nil {
		os.Exit(1)
	}

	flag.StringVar(&args.notation, "notation", "", "notation syntax name")
	flag.StringVar(&args.input, "in", "", "input file name")
	flag.StringVar(&args.data, "data", "", "input text")
	flag.StringVar(&args.output, "out", "chess2pic", "output file name (without extension)")
	flag.StringVar(&args.from, "from", "white", "from which player's perspective (\"white\" or \"black\") to draw")

	flag.BoolVar(&chess2pic.DEBUG, "debug", false, "enable debug output")
}

func main() {
	flag.Parse()

	if args.notation == "" {
		chess2pic.Fatalf("--notation is required")
	}

	var from chess.PieceColor
	switch args.from {
	case "white":
		from = chess.White
	case "black":
		from = chess.Black
	default:
		chess2pic.Fatalf("invalid --from value: %q", args.from)
	}

	var in io.Reader
	if args.input != "" {
		f, err := os.Open(args.input)
		if err != nil {
			chess2pic.Fatalf("error opening %q: %s", args.input, err)
		}
		defer f.Close()
		in = bufio.NewReader(f)
	} else {
		in = strings.NewReader(args.data)
	}

	var err error
	switch args.notation {
	case "fen":
		err = chess2pic.HandleFEN(in, args.output, pic.DefaultCollection, from)
	case "pgn":
		err = chess2pic.HandlePGN(in, args.output, pic.DefaultCollection, from)
	default:
		err = fmt.Errorf("unknown notation: %q", args.notation)
	}
	if err != nil {
		chess2pic.Fatalf(err.Error())
	}
}
