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

func TestSquareString(t *testing.T) {
	tcs := []struct {
		file int
		rank int
		s    string
	}{
		{0, 0, "a1"},
		{4, 3, "e4"},
		{7, 7, "h8"},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("test case (%d, %d)", tc.file, tc.rank), func(tt *testing.T) {
			got := MustNewSquare(tc.file, tc.rank).String()
			if got != tc.s {
				tt.Errorf("want %q, got %q", tc.s, got)
			}
		})
	}
}

func TestMustNewSquareFromString(t *testing.T) {
	tcs := []struct {
		s    string
		file int
		rank int
	}{
		{"a1", 0, 0},
		{"A1", 0, 0},
		{"g5", 6, 4},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("test case %q", tc.s), func(tt *testing.T) {
			sq := MustNewSquareFromString(tc.s)
			if sq.file != tc.file {
				tt.Errorf("wrong file: want %d, got %d", tc.file, sq.file)
			}
			if sq.rank != tc.rank {
				tt.Errorf("wrong rank: want %d, got %d", tc.rank, sq.rank)
			}
		})
	}
}

func TestOnDiag(t *testing.T) {
	tcs := []struct {
		a    string
		b    string
		want bool
	}{
		{"e4", "c6", true},
		{"e4", "g6", true},
		{"e4", "c2", true},
		{"e4", "g2", true},
		{"e4", "c5", false},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%s-%s", tc.a, tc.b), func(tt *testing.T) {
			a, err := NewSquareFromString(tc.a)
			if err != nil {
				panic(err)
			}
			b, err := NewSquareFromString(tc.b)
			if err != nil {
				panic(err)
			}
			if got := OnDiag(a, b); got != tc.want {
				tt.Errorf("(a, b): want %t, got %t", tc.want, got)
			}
			if got := OnDiag(b, a); got != tc.want {
				tt.Errorf("(b, a): want %t, got %t", tc.want, got)
			}
		})
	}
}

func TestApply(t *testing.T) {

	tcs := []struct {
		name  string
		start []squarePiece
		mov   move
		end   []squarePiece
	}{
		{
			name:  "simple move",
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
			name:  "promotion",
			start: []squarePiece{{"e7", Piece{Pawn, White}}},
			mov:   move{from: "e7", to: "e8", pr: Piece{Queen, White}},
			end:   []squarePiece{{"e8", Piece{Queen, White}}},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			start := getPosition(tc.start)
			mov := getMove(tc.mov)
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
