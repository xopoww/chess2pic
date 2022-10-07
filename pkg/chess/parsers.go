package chess

import "io"

// PositionParser parses UTF-8 representation of chess position
type PositionParser interface {
	Parse(r io.RuneReader) (Position, error)
}

// MoveParser parses UTF-8 representation of series of chess moves
type MoveParser interface {
	Parse(start Position, r io.RuneReader) ([]Move, error)
}
