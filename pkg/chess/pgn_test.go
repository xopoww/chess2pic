package chess

import (
	"strings"
	"testing"
)

func TestPgnParse(t *testing.T) {

	tcs := []struct {
		name     string
		notation string
		want     pgnResult
	}{
		{
			name:     "from starting position",
			notation: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			},
		},
		{
			name:     "with tags",
			notation: "[Foo \"bar\"]\n[Baz \"quux\"]\n\n1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
				tags: map[string]string{
					"Foo": "bar",
					"Baz": "quux",
				},
			},
		},
		{
			name:     "from position",
			notation: "[FEN \"k7/1p6/8/8/8/8/6P1/7K w - - 0 1\"]\n\n1. g4 b5 2. g5 b4",
			want: pgnResult{
				start: "k7/1p6/8/8/8/8/6P1/7K",
				movs:  "1. g4 b5 2. g5 b4",
				tags: map[string]string{
					"FEN": "k7/1p6/8/8/8/8/6P1/7K w - - 0 1",
				},
			},
		},
		{
			name:     "tags with escape sequences",
			notation: "[Foo \"ba\\\"r\"]\n[Baz \"qu\\\\ux\"]\n\n1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
				tags: map[string]string{
					"Foo": "ba\"r",
					"Baz": "qu\\ux",
				},
			},
		},
		{
			name: "with comment",
			notation: "1. e4 e5 2. Nf3 Nf6 { Something really smart } 3. Nxe5 Nc6 4. Nxc6 dxc6",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			},
		},
		{
			name: "with game result (win)",
			notation: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6 1-0",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			},
		},
		{
			name: "with game result (draw)",
			notation: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6 1/2-1/2",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			},
		},
		{
			name: "with game result (ongoing)",
			notation: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6 *",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 dxc6",
			},
		},
		{
			name: "with game result (ongoing on black's turn)",
			notation: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6 *",
			want: pgnResult{
				movs: "1. e4 e5 2. Nf3 Nf6 3. Nxe5 Nc6 4. Nxc6",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			want := getPgnResult(tc.want)
			r := strings.NewReader(tc.notation)
			got, err := ParsePGN(r)
			if err != nil {
				tt.Fatal(err)
			}
			if !positionEqual(want.Start, got.Start) {
				tt.Fatalf("want:\n%s\ngot:\n%s\n", want.Start, got.Start)
			}
			assertMoves(tt, want.Moves, got.Moves)
			for k, v := range want.Tags {
				gv, exists := got.Tags[k]
				if !exists {
					tt.Errorf("missing tag %q", k)
					continue
				}
				if v != gv {
					tt.Errorf("tag %q: want %q, got %q", k, v, gv)
				}
				delete(got.Tags, k)
			}
			if len(got.Tags) > 0 {
				tt.Errorf("extra tags: %#v", got.Tags)
			}
		})
	}

}
