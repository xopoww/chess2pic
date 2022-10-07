package chess

import (
	"strings"
	"testing"
)

func TestFenParserParse(t *testing.T) {
	tcs := []struct {
		name     string
		notation string
		want     Position
		wantErr  error
	}{
		{
			name:     "starting position",
			notation: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			want:     StartingPosition(),
		},
		{
			name:     "starting position with metainfo",
			notation: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want:     StartingPosition(),
		},
		{
			name:     "too long rank",
			notation: "rnbqkbnrP/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			wantErr:  ErrTooLongRank,
		},
		{
			name:     "too short rank",
			notation: "rnbqkbn/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			wantErr:  ErrTooShortRank,
		},
		{
			name:     "too many ranks",
			notation: "rnbqkbnr/pppppppp/8/8/8/8/8/PPPPPPPP/RNBQKBNR",
			wantErr:  ErrTooManyRanks,
		},
		{
			name:     "too few ranks",
			notation: "rnbqkbnr/pppppppp/8/8/8/PPPPPPPP/RNBQKBNR",
			wantErr:  ErrTooFewRanks,
		},
		{
			name:     "invalid character",
			notation: "rnbqkbnr/pppppppp/8/8/8/8/PPP🌚PPPP/RNBQKBNR",
			wantErr:  InvalidRuneError{At: 29, Rune: '🌚'},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			r := strings.NewReader(tc.notation)
			pos, err := FEN().Parse(r)
			if tc.wantErr != err {
				tt.Fatalf("want error: %v, got error: %v", tc.wantErr, err)
			}
			if err == nil && !positionEqual(tc.want, pos) {
				tt.Fatalf("\nwant:\n%s\ngot:\n%s\n", tc.want, pos)
			}
		})
	}
}
