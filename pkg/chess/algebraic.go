package chess

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	number = iota
	white
	black
)


type IllegalMoveError struct {
	FullmoveIndex int
	Color		  PieceColor
	Notation	  string
}

func (err IllegalMoveError) Error() string {
	return fmt.Sprintf("illegal move #%d for %s: %q", err.FullmoveIndex + 1, err.Color.Name(), err.Notation)
}


type InvalidSyntaxError struct {
	At int
	Reason string
}

func (err InvalidSyntaxError) Error() string {
	return fmt.Sprintf("invalid SAN syntax at %d: %s", err.At, err.Reason)
}


type algParser struct {
	state int
	lastNum int
	nread int
	pos  Position
	movs []Move
}

func Algebraic() MoveParser {
	return &algParser{}
}

func (ap *algParser) addMove(mov Move) {
	ap.movs = append(ap.movs, mov)
	ap.pos = Apply(ap.pos, mov)
}

func (ap *algParser) handleNumber(cs []rune) error {
	s := string(cs)

	nextState := white
	if strings.HasSuffix(s, "...") {
		if ap.lastNum < 0 {
			s = strings.TrimSuffix(s, "...")
			nextState = black
		} else {
			return InvalidSyntaxError{At: ap.nread - len(s), Reason: "unexpected \"...\" in this context"}
		}
	} else if strings.HasSuffix(s, ".") {
		s = strings.TrimSuffix(s, ".")
	} else {
		return InvalidSyntaxError{At: ap.nread - len(s), Reason: fmt.Sprintf("invalid move number notation: %q", s)}
	}

	num, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return InvalidSyntaxError{At: ap.nread - len(s), Reason: fmt.Sprintf("invalid move number notation: %q", s)}
	}
	if ap.lastNum > 0 && ap.lastNum + 1 != int(num) {
		return InvalidSyntaxError{At: ap.nread - len(s), Reason: fmt.Sprintf("expected move #%d, got #%d", ap.lastNum + 1, num)}
	}

	ap.lastNum = int(num)
	ap.state = nextState

	return nil
}

func pieceByLetter(letter rune) PieceKind {
	switch letter {
	case 'R':
		return Rook
	case 'N':
		return Knight
	case 'B':
		return Bishop
	case 'Q':
		return Queen
	case 'K':
		return King
	default:
		return Pawn
	}
}

func findSources(pos Position, p Piece, destination Square, capture bool) []Square {
	sources := make([]Square, 0)

	// find king for pinned pieces checks
	var king Square
	findKing:
	for file := range pos {
		for rank := range pos[file] {
			sq := MustNewSquare(file, rank)
			pp := pos.Get(sq)
			if pp.Kind == King && pp.Color == p.Color {
				king = sq
				break findKing
			}
		}
	}

	isPinned := func(sq Square) bool {
		if p.Kind == King {
			return false
		}

		// check lateral pin (same file)
		if king.file == sq.file && king.file != destination.file {
			var dr int
			if king.rank < sq.rank {
				dr = 1
			} else {
				dr = -1
			}
			for rank := sq.rank + dr; rank >= 0 && rank <= 7; rank += dr {
				pp := pos.Get(MustNewSquare(sq.file, rank))
				if pp.Color != p.Color && (pp.Kind == Rook || pp.Kind == Queen) {
					return true
				} else if pp.Kind != None {
					break
				}
			}
		// check lateral pin (same rank)
		} else if king.rank == sq.rank && king.rank != destination.rank {
			var df int
			if king.file < sq.file {
				df = 1
			} else {
				df = -1
			}
			for file := sq.file + df; file >= 0 && file <= 7; file += df {
				pp := pos.Get(MustNewSquare(file, sq.rank))
				if pp.Color != p.Color && (pp.Kind == Rook || pp.Kind == Queen) {
					return true
				} else if pp.Kind != None {
					break
				}
			}
		// check diagonal pin
		} else if OnDiag(king, sq) {
			var (
				dr int
				df int
			)
			if king.rank < sq.rank {
				dr = 1
			} else {
				dr = -1
			}
			if king.file < sq.file {
				df = 1
			} else {
				df = -1
			}
			// check that the piece leaves the diagonal
			leftDiag := !OnDiag(king, destination) || ((king.rank - destination.rank) * (king.file - destination.file)) * (dr * df) < 0
			if leftDiag {
				for d := 1; ; d++ {
					ssq, err := NewSquare(sq.file + d * df, sq.rank + d * dr)
					if err != nil {
						break
					}
					pp := pos.Get(ssq)
					if pp.Color != p.Color && (pp.Kind == Bishop || pp.Kind == Queen) {
						return true
					} else if pp.Kind != None {
						break
					}
				}
			}
		}
		return false
	}
	
	addSource := func(file, rank int) {
		sq, err := NewSquare(file, rank)
		if err != nil {
			return
		}

		if pos.Get(sq) != p {
			return
		}

		sources = append(sources, sq)
	}

	addDiagonal := func() {
		for df := -1; df <= 1; df += 2 {
			for dr := -1; dr <= 1; dr += 2 {
				d := 0
				for {
					d++
					source, err := NewSquare(destination.file + df * d, destination.rank + dr * d)
					if err != nil {
						break
					}
					pp := pos.Get(source)
					if p == pp {
						sources = append(sources, source)
						break
					} else if pp.Kind != None {
						break
					}
				}
			}
		}
	}

	addLateral := func() {
		for df := -1; df <= 1; df += 2 {
			d := 0
			for {
				d++
				source, err := NewSquare(destination.file + df * d, destination.rank)
				if err != nil {
					break
				}
				pp := pos.Get(source)
				if p == pp {
					sources = append(sources, source)
					break
				} else if pp.Kind != None {
					break
				}
			}
		}
		for dr := -1; dr <= 1; dr += 2 {
			d := 0
			for {
				d++
				source, err := NewSquare(destination.file, destination.rank + dr * d)
				if err != nil {
					break
				}
				pp := pos.Get(source)
				if p == pp {
					sources = append(sources, source)
					break
				} else if pp.Kind != None {
					break
				}
			}
		}
	}

	switch p.Kind {
	case Pawn:
		if capture {
			var rank int
			if p.Color == White {
				rank = destination.rank - 1
			} else {
				rank = destination.rank + 1
			}
			for df := -1; df <= 1; df += 2 {
				addSource(destination.file + df, rank)
			}
		} else {
			var dRank int
			if p.Color == White {
				dRank = -1
			} else {
				dRank = 1
			}
			addSource(destination.file, destination.rank + dRank)
			if p.Color == White && destination.rank == 3 || p.Color == Black && destination.rank == 4 {
				addSource(destination.file, destination.rank + dRank * 2)
			}
		}
	case Knight:
		addSource(destination.file - 1, destination.rank - 2)
		addSource(destination.file + 1, destination.rank - 2)
		addSource(destination.file - 1, destination.rank + 2)
		addSource(destination.file + 1, destination.rank + 2)
		addSource(destination.file - 2, destination.rank - 1)
		addSource(destination.file + 2, destination.rank - 1)
		addSource(destination.file - 2, destination.rank + 1)
		addSource(destination.file + 2, destination.rank + 1)
	case Bishop:
		addDiagonal()
	case Rook:
		addLateral()
	case Queen:
		addDiagonal()
		addLateral()
	case King:
		for df := -1; df <= 1; df++ {
			for dr := -1; dr <= 1; dr++ {
				if df == 0 && dr == 8 {
					continue
				}
				addSource(destination.file + df, destination.rank + dr)
			}
		}
	}

	filteredSources := make([]Square, 0)
	for _, sq := range sources {
		if !isPinned(sq) {
			filteredSources = append(filteredSources, sq)
		}
	}
	return filteredSources
}

func (ap *algParser) handleMove(cs []rune) error {
	s := string(cs)

	p := Piece{}
	if ap.state == white {
		p.Color = White
		ap.state = black
	} else {
		p.Color = Black
		ap.state = number
	}

	illegal := IllegalMoveError{FullmoveIndex: ap.lastNum - 1, Color: p.Color, Notation: s}
	if len(cs) < 2 {
		return illegal
	}

	// handle castling
	if s == "O-O" || s == "O-O-O" {
		p.Kind = King
		
		var rank int
		if p.Color == White {
			rank = 0
		} else {
			rank = 7
		}
		var (
			kFile  int = 4  // file of square with the king
			nkFile int      // file of square for the king (must be none)
			rFile  int      // file of square with the rook
			nrFile int      // file of square for the rook (must be none)
		)
		if s == "O-O" {
			nkFile = 6
			rFile = 7
			nrFile = 5
		} else {
			nkFile = 2
			rFile = 0
			nrFile = 3
		}
		
		checks := []struct{
			file int
			want PieceKind
		}{
			{kFile, King}, {nkFile, None}, {rFile, Rook}, {nrFile, None},
		}
		for _, check := range checks {
			pp := ap.pos.Get(MustNewSquare(check.file, rank))
			if pp.Kind != check.want || (check.want != None && pp.Color != p.Color) {
				return illegal
			}
		}
		ap.addMove(Move{
			From: MustNewSquare(kFile, rank),
			To: MustNewSquare(nkFile, rank),
			Castle: true,
		})
		return nil
	}

	// determine which piece moves
	p.Kind = pieceByLetter(cs[0])
	if p.Kind != Pawn {
		cs = cs[1:]
	}

	mov := Move{}

	// check for promotion
	if cs[len(cs)-2] == '=' {
		if p.Kind != Pawn {
			return illegal
		}
		mov.Promotion = Piece{pieceByLetter(cs[len(cs)-1]), p.Color}
		if mov.Promotion.Kind == Pawn {
			return illegal
		}
		cs = cs[:len(cs)-2]
	}

	// need to check len(cs) again after cuts
	if len(cs) < 2 {
		return illegal
	}

	// get destination square
	to, err := NewSquareFromString(string(cs[len(cs)-2:]))
	if err != nil {
		return illegal
	}
	cs = cs[:len(cs)-2]
	mov.To = to

	// check that promotion happens on the right rank
	if mov.Promotion.Kind != None {
		if (p.Color == White && mov.To.rank != 7) || (p.Color == Black && mov.To.rank != 0) {
			return illegal
		}
	}

	// check for capture (and en passant)
	capture := false
	if len(cs) > 0 && cs[len(cs)-1] == 'x' {
		capture = true
		cs = cs[:len(cs)-1]

		target := ap.pos.Get(mov.To)
		if target.Kind == None {
			if p.Kind == Pawn {
				if p.Color == White {
					target = ap.pos.Get(MustNewSquare(mov.To.file, mov.To.rank - 1))
				} else {
					target = ap.pos.Get(MustNewSquare(mov.To.file, mov.To.rank + 1))
				}
				if target.Kind == Pawn {
					mov.EnPassant = true
				} else {
					return illegal
				}
			} else {
				return illegal
			}
		}
		if target.Color == p.Color {
			return illegal
		}
	} else {
		if ap.pos.Get(mov.To).Kind != None {
			return illegal
		}
	}

	// look for potential source squares
	sources := findSources(ap.pos, p, mov.To, capture)
	// look for source disambiguation hints
	var (
		sFile int = -1
		sRank int = -1
	)
	switch len(cs) {
	case 0:
		break
	case 1:
		if cs[0] >= 'a' && cs[0] <= 'h' {
			sFile = int(cs[0] - 'a')
		} else if cs[0] >= '1' && cs[0] <= '8' {
			sRank = int(cs[0] - '1')
		} else {
			return illegal
		}
	case 2:
		if cs[0] >= 'a' && cs[0] <= 'h' {
			sFile = int(cs[0] - 'a')
		} else {
			return illegal
		}
		if cs[1] >= '1' && cs[1] <= '8' {
			sRank = int(cs[1] - '1')
		} else {
			return illegal
		}
	default:
		return illegal
	}
	// filter sources using hints
	filteredSources := make([]Square, 0)
	for _, source := range sources {
		if sFile >= 0 && source.file != sFile {
			continue
		}
		if sRank >= 0 && source.rank != sRank {
			continue
		}
		filteredSources = append(filteredSources, source)
	}
	// must be exactly one source square
	if len(filteredSources) != 1 {
		return illegal
	}
	mov.From = filteredSources[0]

	ap.addMove(mov)
	return nil
}

func (ap *algParser) handle(cs []rune) error {
	switch ap.state {
	case number:
		return ap.handleNumber(cs)
	case white: fallthrough
	case black:
		return ap.handleMove(cs)
	}
	panic("unknown parser state")
}

func (ap *algParser) Parse(start Position, r io.RuneReader) ([]Move, error) {
	ap.pos = start
	ap.movs = make([]Move, 0)

	ap.nread = -1
	ap.state = number
	ap.lastNum = -1

	var cs []rune

	for {
		c, _, err := r.ReadRune()
		// handle EOF
		if errors.Is(err, io.EOF) {
			// handle leftover string
			if len(cs) > 0 {
				if err := ap.handle(cs); err != nil {
					return nil, err
				}
			}
			return ap.movs, nil
		}
		ap.nread++

		// handle whitespace
		if unicode.IsSpace(c) {
			// if there is a string collected - handle it
			if len(cs) > 0 {
				if err := ap.handle(cs); err != nil {
					return nil, err
				}
				cs = cs[:0]
			}
			continue
		}

		// handle a character according to the state
		allowedRunes := map[int]string{
			number: "1234567890.",
			white:  "RNBQK"+"abcdefgh"+"12345678"+"x+#="+"O-",
			black:  "RNBQK"+"abcdefgh"+"12345678"+"x+#="+"O-",
		}

		if strings.ContainsRune(allowedRunes[ap.state], c) {
			cs = append(cs, c)
		} else {
			return nil, InvalidSyntaxError{At: ap.nread, Reason: fmt.Sprintf("unexpected %U in this context", c)}
		}
	}
}