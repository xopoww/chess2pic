package chess

import (
	"fmt"
	"testing"
)

func TestStartingPosition(t *testing.T) {
	pos := StartingPosition()
	const wantString = "brbkbbbqbKbbbkbr\nbpbpbpbpbpbpbpbp\n................\n................\n................\n................\nwpwpwpwpwpwpwpwp\nwrwkwbwqwKwbwkwr\n"
	if pos.String() != wantString {
		t.Errorf("Want:\n%s\nGot:\n%s", wantString, pos.String())
	}
}

func TestNewSquareFromString(t *testing.T) {
	tcs := []struct{
		s string
		file int
		rank int
	}{
		{"a1", 0, 0},
		{"A1", 0, 0},
		{"g5", 6, 4},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("test case %q", tc.s), func(tt *testing.T) {
			sq := NewSquareFromString(tc.s)
			if sq.file != tc.file {
				tt.Errorf("wrong file: want %d, got %d", tc.file, sq.file)
			}
			if sq.rank != tc.rank {
				tt.Errorf("wrong rank: want %d, got %d", tc.rank, sq.rank)
			}
		})
	}
}

func TestApply(t *testing.T) {
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

	tcs :=  []struct {
		name  string
		start []squarePiece
		mov   move	
		end   []squarePiece
	}{
		{
			name: "simple move",
			start: []squarePiece{{"e2", Piece{Pawn, White}}},
			mov:   move{from: "e2", to: "e4"},
			end:   []squarePiece{{"e4", Piece{Pawn, White}}},
		},
		{
			name:  "simple capture",
			start: []squarePiece{{"e2", Piece{Pawn, White}}, {"d3", Piece{Pawn, Black}}},
			mov:   move{from: "e2", to: "d3"},
			end:   []squarePiece{{"d3", Piece{Pawn, White}}},
		},
		{
			name:  "en passant",
			start: []squarePiece{{"e5", Piece{Pawn, White}}, {"f5", Piece{Pawn, Black}}},
			mov:   move{from: "e5", to: "f6", ep: true},
			end:   []squarePiece{{"f6", Piece{Pawn, White}}},
		},
		{
			name:  "short castle",
			start: []squarePiece{{"e1", Piece{King, White}}, {"h1", Piece{Rook, White}}},
			mov:   move{from: "e1", to: "g1", cs: true},
			end:   []squarePiece{{"g1", Piece{King, White}}, {"f1", Piece{Rook, White}}},
		},
		{
			name:  "long castle",
			start: []squarePiece{{"e1", Piece{King, White}}, {"a1", Piece{Rook, White}}},
			mov:   move{from: "e1", to: "c1", cs: true},
			end:   []squarePiece{{"c1", Piece{King, White}}, {"d1", Piece{Rook, White}}},
		},
		{
			name: "promotion",
			start: []squarePiece{{"e7", Piece{Pawn, White}}},
			mov:   move{from: "e7", to: "e8", pr: Piece{Queen, White}},
			end:   []squarePiece{{"e8", Piece{Queen, White}}},
		},
	}

	getPosition := func(sps []squarePiece) Position {
		var pos Position
		for _, sp := range sps {
			pos = pos.Set(NewSquareFromString(sp.sq), sp.p)
		}
		return pos
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			start := getPosition(tc.start)
			mov := Move{
				From: NewSquareFromString(tc.mov.from),
				To:   NewSquareFromString(tc.mov.to),
				EnPassant: tc.mov.ep,
				Castle:    tc.mov.cs,
				Promotion: tc.mov.pr,
			}
			got := Apply(start, mov)
			want := getPosition(tc.end)
			for file := range want {
				for rank := range want[file] {
					wp := want[file][rank]
					gp := got[file][rank]
					if wp != gp {
						tt.Errorf("at [%d][%d]: want %s, got %s", file, rank, wp, gp)
					}
				}
			}
		})
	}
}
