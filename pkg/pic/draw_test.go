package pic

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/xopoww/chess2pic/pkg/chess"
)

type mockColor byte

func (mc mockColor) RGBA() (r uint32, g uint32, b uint32, a uint32) {
	return uint32(mc), 0, 0, 0xffff
}

type mockImage struct {
	size int
	data [64]byte
}

func (mi *mockImage) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		r, _, _, _ := c.RGBA()
		return mockColor(r)
	})
}

func (mi *mockImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, mi.size, mi.size)
}

func (mi *mockImage) At(x int, y int) color.Color {
	return mockColor(mi.data[x+y*mi.size])
}

func (mi *mockImage) Set(x int, y int, c color.Color) {
	mi.data[x+y*mi.size] = byte(mi.ColorModel().Convert(c).(mockColor))
}

type mockCollection struct{}

func (mcol mockCollection) Board(fromPerspective chess.PieceColor) Image {
	data := [64]byte{}
	for i := range data {
		data[i] = 0xFF
	}
	return &mockImage{
		size: 8,
		data: data,
	}
}

func (mcol mockCollection) Piece(p chess.Piece) Image {
	b := byte(p.Kind)
	if p.Color == chess.Black {
		b |= 0x10
	}
	return &mockImage{
		size: 1,
		data: [64]byte{b},
	}
}

func (mcol mockCollection) Canvas() draw.Image {
	return &mockImage{
		size: 8,
	}
}

func TestDrawPosition(t *testing.T) {
	col := mockCollection{}
	pos := chess.StartingPosition()

	for from := chess.White; from <= chess.Black; from++ {
		t.Run(fmt.Sprintf("from %s", from.Name()), func(tt *testing.T) {
			ddst := DrawPosition(col, pos, from)
			dst, ok := ddst.(*mockImage)
			if !ok {
				tt.Fatalf("wrong draw.Image type (got %#v)", ddst)
			}

			var wantDataHex string
			if from == chess.White {
				wantDataHex =
					"1213141516141312" +
						"1111111111111111" +
						"ffffffffffffffff" +
						"ffffffffffffffff" +
						"ffffffffffffffff" +
						"ffffffffffffffff" +
						"0101010101010101" +
						"0203040506040302"
			} else {
				wantDataHex =
					"0203040605040302" +
						"0101010101010101" +
						"ffffffffffffffff" +
						"ffffffffffffffff" +
						"ffffffffffffffff" +
						"ffffffffffffffff" +
						"1111111111111111" +
						"1213141615141312"
			}
			wantData, err := hex.DecodeString(wantDataHex)
			if err != nil {
				panic(err)
			}

			gotData := dst.data
			if len(gotData) != len(wantData) {
				tt.Fatalf("want len = %d, got len = %d", len(wantData), len(gotData))
			}
			nMismathes := 0
			for i := range wantData {
				if wantData[i] != gotData[i] {
					nMismathes++
					tt.Errorf("mismatch at %d: want %02x, got %02x", i, wantData[i], gotData[i])
				}
			}
			if nMismathes > 0 {
				tt.Errorf("total of %d mismatches", nMismathes)
			}
		})
	}
}
