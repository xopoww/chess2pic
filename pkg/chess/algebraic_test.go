package chess

import (
	"strings"
	"testing"
)


func TestAlgParserParse(t *testing.T) {
	const skipErrorTests = false

	tcs := []struct {
		name	 string
		start 	 string	// in FEN
		notation string
		want     []move
		wantErr  bool
	}{
		{
			name: "pawn moves",
			start: "4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3",
			notation: "1. e4 d6 2. c4 d5 3. cxd5 c5 4. dxc6 f5 5. c7 fxe4 6. c8=N",
			want: []move{
				{from: "e2", to: "e4"},
				{from: "d7", to: "d6"},
				{from: "c2", to: "c4"},
				{from: "d6", to: "d5"},
				{from: "c4", to: "d5"},
				{from: "c7", to: "c5"},
				{from: "d5", to: "c6", ep: true},
				{from: "f7", to: "f5"},
				{from: "c6", to: "c7"},
				{from: "f5", to: "e4"},
				{from: "c7", to: "c8", pr: Piece{Knight, White}},
			},
		},
		{
			name: "piece moves",
			start: "rnbqkb2/8/8/8/8/8/8/RNBQKB2",
			notation: "1. Ra4 Ra5 2. Nd2 Nc6 3. Bb2 Be6 4. Ba6 Ba3 5. Qg4 Qd4 6. Kd1 Kf7 7. Rxa5 Nxa5 8. Bxa3",
			want: []move{
				{from: "a1", to: "a4"},
				{from: "a8", to: "a5"},
				{from: "b1", to: "d2"},
				{from: "b8", to: "c6"},
				{from: "c1", to: "b2"},
				{from: "c8", to: "e6"},
				{from: "f1", to: "a6"},
				{from: "f8", to: "a3"},
				{from: "d1", to: "g4"},
				{from: "d8", to: "d4"},
				{from: "e1", to: "d1"},
				{from: "e8", to: "f7"},
				{from: "a4", to: "a5"},
				{from: "c6", to: "a5"},
				{from: "b2", to: "a3"},
			},
		},
		{
			name: "castling",
			start: "r3k2r/8/8/8/8/8/8/R3K2R",
			notation: "1. O-O O-O-O",
			want: []move{
				{from: "e1", to: "g1", cs: true},
				{from: "e8", to: "c8", cs: true},
			},
		},
		{
			name: "disambiguation",
			start: "r5k1/8/r7/2b3b1/2N5/1N6/R6R/7K",
			notation: "1. Rae2 R8a7 2. Ncd2 Bce7",
			want: []move{
				{from: "a2", to: "e2"},
				{from: "a8", to: "a7"},
				{from: "c4", to: "d2"},
				{from: "c5", to: "e7"},
			},
		},
		{
			name: "disambiguation lvl 2",
			start: "kn6/pp6/8/8/8/5Q2/8/K2Q1Q2",
			notation: "1. Qf1d3",
			want: []move{
				{from: "f1", to: "d3"},
			},
		},
		{
			name: "disambiguation lvl 3",	// tricky edge case with pinned piece
			start: "k7/8/8/8/q7/R6R/K7/8",
			notation: "1. Rd3",
			want: []move{
				{from: "h3", to: "d3"},
			},
		},
		{
			name: "black moves first",
			start: "k7/p7/8/8/8/8/P7/K7",
			notation: "1... a5",
			want: []move{
				{from: "a7", to: "a5"},
			},
		},
		{
			name: "check",
			start: "1k6/7Q/2K5/8/8/8/8/8",
			notation: "1. Qh8+",
			want: []move{
				{from: "h7", to: "h8"},
			},
		},
		{
			name: "checkmate",
			start: "1k6/7Q/2K5/8/8/8/8/8",
			notation: "1. Qb7#",
			want: []move{
				{from: "h7", to: "b7"},
			},
		},
		{
			name: "promotion checkmate",
			start: "1k6/7P/1K6/8/8/8/8/8",
			notation: "1. h8=Q#",
			want: []move{
				{from: "h7", to: "h8", pr: Piece{Queen, White}},
			},
		},
		{
			name: "capture with promotion checkmate",
			start: "1k4n1/5P1P/1K6/8/8/8/8/8",
			notation: "1. hxg8=Q#",
			want: []move{
				{from: "h7", to: "g8", pr: Piece{Queen, White}},
			},
		},

		{
			name: "illegal pawn move",
			start: "k7/8/8/8/8/8/P7/NRBB1Q1K",
			notation: "1. e4",
			wantErr: true,
		},
		{
			name: "illegal rook move",
			start: "k7/8/8/8/8/8/P7/NRBB1Q1K",
			notation: "1. Re4",
			wantErr: true,
		},
		{
			name: "illegal knight move",
			start: "k7/8/8/8/8/8/P7/NRBB1Q1K",
			notation: "1. Ne4",
			wantErr: true,
		},
		{
			name: "illegal bishop move",
			start: "k7/8/8/8/8/8/P7/NRBB1Q1K",
			notation: "1. Be4",
			wantErr: true,
		},
		{
			name: "illegal queen move",
			start: "k7/8/8/8/8/8/P7/NRBB1Q1K",
			notation: "1. Qe4",
			wantErr: true,
		},
		{
			name: "illegal king move",
			start: "k7/8/8/8/8/8/P7/NRBB1Q1K",
			notation: "1. Ke4",
			wantErr: true,
		},
		{
			name: "wrong rank promotion",
			start: "k7/8/8/P7/8/8/8/7K",
			notation: "1. a6=N",
			wantErr: true,
		},
		{
			name: "wrong rank long pawn move",
			start: "k7/8/8/P7/8/8/8/7K",
			notation: "1. a7",
			wantErr: true,
		},

		{
			name: "lateral pin (same file)",
			start: "k6q/8/8/8/8/8/7R/7K",
			notation: "1. Rb2",
			wantErr: true,
		},
		{
			name: "lateral pin (same file) (legal)",
			start: "k6q/8/8/8/8/8/7R/7K",
			notation: "1. Rh4",
			want: []move{
				{from: "h2", to: "h4"},
			},
		},
		{
			name: "lateral pin (same rank)",
			start: "k7/8/8/8/8/8/q5RK/8",
			notation: "1. Rg5",
			wantErr: true,
		},
		{
			name: "lateral pin (same rank) (legal)",
			start: "k7/8/8/8/8/8/q5RK/8",
			notation: "1. Rb2",
			want: []move{
				{from: "g2", to: "b2"},
			},
		},
		{
			name: "diagonal pin",
			start: "1k6/8/q7/8/2Q5/8/4K3/8",
			notation: "1. Qe4",
			wantErr: true,
		},
		{
			name: "diagonal pin (other diag)",
			start: "1k6/8/q7/8/2Q5/8/4K3/8",
			notation: "1. Qg4",
			wantErr: true,
		},
		{
			name: "diagonal pin (legal)",
			start: "1k6/8/q7/8/2Q5/8/4K3/8",
			notation: "1. Qd3",
			want: []move{
				{from: "c4", to: "d3"},
			},
		},
	}


	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			if skipErrorTests && tc.wantErr {
				tt.SkipNow()
			}

			start, err := FEN().Parse(strings.NewReader(tc.start))
			if err != nil {
				panic(err)
			}
			r := strings.NewReader(tc.notation)
			want := make([]Move, 0, len(tc.want))
			for _, mov := range tc.want {
				want = append(want, getMove(mov))
			}

			got, err := Algebraic().Parse(start, r)
			if tc.wantErr != (err != nil) {
				tt.Fatalf("want error: %t, got error: %v", tc.wantErr, err)
			}
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
		})
	}
}