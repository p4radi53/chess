package chess

type CastlingPossibility struct {
	IsWhiteKingsidePossible  bool
	IsBlackKingsidePossible  bool
	IsWhiteQueensidePossible bool
	IsBlackQueensidePossible bool
}

func (c *CastlingPossibility) updateCastlingPossibility(lastMove *Move) {
	if !c.IsBlackKingsidePossible && !c.IsWhiteKingsidePossible && !c.IsWhiteQueensidePossible && !c.IsBlackQueensidePossible {
		return
	}
	movedPiece := lastMove.ColoredPiece

	if movedPiece.Piece != King && movedPiece.Piece != Rook {
		return
	}

	if movedPiece.Piece == King {
		if movedPiece.Color == White {
			c.IsWhiteKingsidePossible = false
			c.IsWhiteQueensidePossible = false
		}
		if movedPiece.Color == Black {
			c.IsBlackKingsidePossible = false
			c.IsBlackQueensidePossible = false
		}
	} else if movedPiece.Piece == Rook {
		lastMoveFromSquare := lastMove.OldSquare
		if movedPiece.Color == White {
			if lastMoveFromSquare == (Square{File: 0, Rank: 0}) {
				c.IsWhiteQueensidePossible = false
			}
			if lastMoveFromSquare == (Square{File: 7, Rank: 0}) {
				c.IsWhiteKingsidePossible = false
			}
		}
		if movedPiece.Color == Black {
			if lastMoveFromSquare == (Square{File: 0, Rank: 7}) {
				c.IsBlackKingsidePossible = false
			}
			if lastMoveFromSquare == (Square{File: 7, Rank: 7}) {
				c.IsBlackQueensidePossible = false
			}
		}
	}
}
