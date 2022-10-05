package chess

import (
	"fmt"
	"strings"
)


type PieceKind int

const (
	None PieceKind = iota
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

func (kind PieceKind) String() string {
	return [...]string{
		" ", "p", "r", "k", "b", "q", "K",
	}[kind]
}

func (kind PieceKind) Name() string {
	return [...]string{
		" ",
		"pawn",
		"rook",
		"knight",
		"bishop",
		"queen",
		"king",
	}[kind]
}


type PieceColor int

const (
	White PieceColor = iota
	Black
)

func (color PieceColor) String() string {
	return [...]string{"w", "b"}[color]
}

func (color PieceColor) Name() string {
	return [...]string{"white", "black"}[color]
}


type Piece struct {
	Kind PieceKind
	Color PieceColor
}

func (p Piece) String() string {
	if p.Kind == None {
		return ".."
	}
	return p.Color.String() + p.Kind.String()
}


type Square struct {
	file int
	rank int
}

func (sq Square) String() string {
	return fmt.Sprintf("%c%d", 'a' + sq.file, sq.rank + 1)
}

// NewSquare creates Square from file and rank coordinates.
// Both coordinates are represented as integers from 0 ("A" file or 1st rank) to 7 ("H" file or 8th rank). 
// If either of coordinates falls out of this range, a zero-value ("A1") is returned.
func NewSquare(file int, rank int) Square {
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return Square{}
	}
	return Square{
		file: file,
		rank: rank,
	}
}

// NewSquareFromString creates square from its string representation (e.g. "G5").
// Both lowercase and uppercase file letters are allowed. If s is not a valid square string, a zero-value ("A1") is returned.
func NewSquareFromString(s string) Square {
	if len(s) != 2 {
		return Square{}
	}
	s = strings.ToLower(s)
	return NewSquare(int(s[0] - 'a'), int(s[1] - '1'))
}


type Position [8][8]Piece

func (pos Position) Get(s Square) Piece {
	return pos[s.file][s.rank]
}

func (pos Position) Set(s Square, p Piece) Position {
	pos[s.file][s.rank] = p
	return pos
}

func (pos Position) String() string {
	bldr := strings.Builder{}
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			bldr.WriteString(pos[file][rank].String())
		}
		bldr.WriteRune('\n')
	}
	return bldr.String()
}


func StartingPosition() Position {
	var pos Position
	
	for file := 0; file < 8; file++ {
		pos[file][1] = Piece{Kind: Pawn, Color: White}
		pos[file][6] = Piece{Kind: Pawn, Color: Black}
	}
	
	for file, kind := range []PieceKind{Rook, Knight, Bishop} {
		pos[file][0]	= Piece{Kind: kind, Color: White}
		pos[7-file][0]	= Piece{Kind: kind, Color: White}
		pos[file][7]	= Piece{Kind: kind, Color: Black}
		pos[7-file][7]	= Piece{Kind: kind, Color: Black}
	}

	pos[3][0] = Piece{Kind: Queen, Color: White}
	pos[3][7] = Piece{Kind: Queen, Color: Black}
	pos[4][0] = Piece{Kind: King, Color: White}
	pos[4][7] = Piece{Kind: King, Color: Black}
	
	return pos
}


type Move struct {
	From Square
	To   Square

	EnPassant bool
	Castle	  bool
	Promotion Piece
}

func (mov Move) String() string {
	s := fmt.Sprintf("%s -> %s", mov.From, mov.To)
	if mov.EnPassant {
		return s + " (e.p.)"
	}
	if mov.Castle {
		return s + " (castle)"
	}
	if mov.Promotion.Kind != None {
		return s + fmt.Sprintf(" (=%s)", mov.Promotion)
	}
	return s
}

func Apply(pos Position, mov Move) Position {
	p := pos.Get(mov.From)
	pos = pos.Set(mov.From, Piece{})

	if mov.Promotion.Kind != None {
		p = mov.Promotion
	}
	pos = pos.Set(mov.To, p)

	if mov.Castle {
		rookMove := Move{}
		rookMove.From.rank = mov.From.rank
		rookMove.To.rank = mov.From.rank

		if mov.To.file == 6 {
			// O-O
			rookMove.From.file = 7
			rookMove.To.file = 5
		} else {
			// O-O-O
			rookMove.From.file = 0
			rookMove.To.file = 3
		}

		return Apply(pos, rookMove)
	}

	if mov.EnPassant {
		captured := mov.To
		if p.Color == White {
			captured.rank--
		} else {
			captured.rank++
		}
		pos = pos.Set(captured, Piece{})
	}

	return pos
}