package chess

import "strings"


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


type PieceColor int

const (
	White PieceColor = iota
	Black
)

func (color PieceColor) String() string {
	return [...]string{"w", "b"}[color]
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

// NewSquare creates Square from file and rank coordinates.
// Both coordinates are represented as integers from 0 ("A" file or 1st rank) to 7 ("H" file or 8th rank). 
// If either of coordinates falls out of this range, a zero-value ("A1") is retuned.
func NewSquare(file int, rank int) Square {
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return Square{}
	}
	return Square{
		file: file,
		rank: rank,
	}
}


type Position [8][8]Piece

func (pos Position) Get(s Square) Piece {
	return pos[s.file][s.rank]
}

func (pos Position) Set(s Square, p Piece) {
	pos[s.file][s.rank] = p
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