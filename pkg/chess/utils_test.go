package chess

type squarePiece struct{
	sq string
	p  Piece
}

type move struct {
	from string
	to 	 string
	ep 	 bool
	cs 	 bool
	pr 	 Piece
}

func getPosition(sps []squarePiece) Position {
	var pos Position
	for _, sp := range sps {
		pos = pos.Set(MustNewSquareFromString(sp.sq), sp.p)
	}
	return pos
}

func getMove(mov move) Move {
	from, err := NewSquareFromString(mov.from)
	if err != nil {
		panic(err)
	}
	to, err := NewSquareFromString(mov.to)
	if err != nil {
		panic(err)
	}
	return Move{
		From: from,
		To:   to,
		EnPassant: mov.ep,
		Castle:    mov.cs,
		Promotion: mov.pr,
	}
}

func positionEqual(a, b Position) bool {
	for file := range a {
		for rank := range a[file] {
			ap := a[file][rank]
			bp := b[file][rank]
			if ap == bp || (ap.Kind == None && bp.Kind == None) {
				continue
			} else {
				return false
			}
		}
	}
	return true
}