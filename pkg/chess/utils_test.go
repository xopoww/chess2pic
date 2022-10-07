package chess

import (
	"strings"
	"testing"
)

type squarePiece struct{
	sq string
	p  Piece
}

type move struct {
	from string
	to 	 string
	ep 	 bool
	cs 	 bool
	pr 	 Piece
}

func getPosition(sps []squarePiece) Position {
	var pos Position
	for _, sp := range sps {
		pos = pos.Set(MustNewSquareFromString(sp.sq), sp.p)
	}
	return pos
}

func getMove(mov move) Move {
	from, err := NewSquareFromString(mov.from)
	if err != nil {
		panic(err)
	}
	to, err := NewSquareFromString(mov.to)
	if err != nil {
		panic(err)
	}
	return Move{
		From: from,
		To:   to,
		EnPassant: mov.ep,
		Castle:    mov.cs,
		Promotion: mov.pr,
	}
}

func positionEqual(a, b Position) bool {
	for file := range a {
		for rank := range a[file] {
			ap := a[file][rank]
			bp := b[file][rank]
			if ap == bp || (ap.Kind == None && bp.Kind == None) {
				continue
			} else {
				return false
			}
		}
	}
	return true
}

type pgnResult struct {
	start string
	movs  string
	tags  map[string]string
}

func getPgnResult(res pgnResult) PGNResult {
	Res := PGNResult{}
	if res.start == "" {
		Res.Start = StartingPosition()
	} else {
		start, err := FEN().Parse(strings.NewReader(res.start))
		if err != nil {
			panic(err)
		}
		Res.Start = start
	}

	movs, err := Algebraic().Parse(Res.Start, strings.NewReader(res.movs))
	if err != nil {
		panic(err)
	}
	Res.Moves = movs

	Res.Tags = res.tags
	return Res
}

func assertMoves(tt *testing.T, want, got []Move) {
	for i := 0; i < len(want) && i < len(got); i++ {
		if want[i] != got[i] {
			tt.Errorf("at [%d]: want %q, got %q", i, want[i], got[i])
		}
	}
	for i := len(want); i < len(got); i++ {
		tt.Errorf("extra move: %q", got[i])
	}
	for i := len(got); i < len(want); i++ {
		tt.Errorf("missing move: %q", want[i])
	}
}