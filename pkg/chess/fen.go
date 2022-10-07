package chess

import (
	"errors"
	"fmt"
	"io"
)

type fenParser struct {
}

func FEN() PositionParser {
	return fenParser{}
}

var (
	ErrTooManyRanks = errors.New("too many ranks")
	ErrTooFewRanks  = errors.New("too few ranks")
	ErrTooLongRank  = errors.New("too long rank")
	ErrTooShortRank = errors.New("too short rank")
)

type InvalidRuneError struct {
	At   int
	Rune rune
}

func (err InvalidRuneError) Error() string {
	return fmt.Sprintf("invalid unicode character %U at position %d", err.Rune, err.At)
}

func (fp fenParser) Parse(r io.RuneReader) (Position, error) {
	pos := Position{}

	at := -1
	file := 0
	rank := 7
	for {
		c, _, err := r.ReadRune()
		// handle error
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return pos, fmt.Errorf("ReadRune: %w", err)
		}
		at++

		// throw away metadata
		if c == ' ' {
			break
		}

		// check that the ranks are not filled
		if rank < 0 {
			return pos, ErrTooManyRanks
		}

		// move to new rank
		if c == '/' {
			if file == 8 {
				file = 0
				rank--
				continue
			} else {
				return pos, ErrTooShortRank
			}
		}

		// handle empty squares
		if c >= '1' && c <= '8' {
			n := int(c - '0')
			if file+n-1 > 7 {
				return pos, ErrTooLongRank
			}
			file += n
			continue
		}

		// handle piece
		pieces := map[rune]Piece{
			'p': {Pawn, Black},
			'n': {Knight, Black},
			'b': {Bishop, Black},
			'r': {Rook, Black},
			'q': {Queen, Black},
			'k': {King, Black},

			'P': {Pawn, White},
			'N': {Knight, White},
			'B': {Bishop, White},
			'R': {Rook, White},
			'Q': {Queen, White},
			'K': {King, White},
		}
		if p, exists := pieces[c]; exists {
			if file > 7 {
				return pos, ErrTooLongRank
			}
			pos = pos.Set(MustNewSquare(file, rank), p)
			file += 1
			continue
		}

		// bad rune
		return pos, InvalidRuneError{At: at, Rune: c}
	}

	if rank == 0 && file == 8 {
		return pos, nil
	} else {
		if rank > 0 {
			return pos, ErrTooFewRanks
		} else {
			return pos, ErrTooShortRank
		}
	}
}
