package chess

import (
	"testing"
)

func TestStartingPosition(t *testing.T) {
	pos := StartingPosition()
	const wantString = "brbkbbbqbKbbbkbr\nbpbpbpbpbpbpbpbp\n................\n................\n................\n................\nwpwpwpwpwpwpwpwp\nwrwkwbwqwKwbwkwr\n"
	if pos.String() != wantString {
		t.Errorf("Want:\n%s\nGot:\n%s", wantString, pos.String())
	}
}