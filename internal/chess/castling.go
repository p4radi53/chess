package chess

type CastlingPossibility struct {
	IsWhiteKingsidePossible  bool
	IsBlackKingsidePossible  bool
	IsWhiteQueensidePossible bool
	IsBlackQueensidePossible bool
}

func castlingMoves(castlingPossibility CastlingPossibility, board Board, color Color) []Square {
	var moves []Square
	opponent := color.Opponent()
	rank := 0
	if color == Black {
		rank = 7
	}

	if color == White {
		if castlingPossibility.IsWhiteKingsidePossible {
			if board.IsCellEmpty(5, rank) && board.IsCellEmpty(6, rank) &&
				!board.IsSquareUnderAttack(Square{5, rank}, opponent) &&
				!board.IsSquareUnderAttack(Square{6, rank}, opponent) {
				moves = append(moves, Square{6, rank})
			}
		}
		if castlingPossibility.IsWhiteQueensidePossible {
			if board.IsCellEmpty(3, rank) && board.IsCellEmpty(2, rank) && board.IsCellEmpty(1, rank) &&
				!board.IsSquareUnderAttack(Square{3, rank}, opponent) &&
				!board.IsSquareUnderAttack(Square{2, rank}, opponent) {
				moves = append(moves, Square{2, rank})
			}
		}
	} else {
		if castlingPossibility.IsBlackKingsidePossible {
			if board.IsCellEmpty(5, rank) && board.IsCellEmpty(6, rank) &&
				!board.IsSquareUnderAttack(Square{5, rank}, opponent) &&
				!board.IsSquareUnderAttack(Square{6, rank}, opponent) {
				moves = append(moves, Square{6, rank})
			}
		}
		if castlingPossibility.IsBlackQueensidePossible {
			if board.IsCellEmpty(3, rank) && board.IsCellEmpty(2, rank) && board.IsCellEmpty(1, rank) &&
				!board.IsSquareUnderAttack(Square{3, rank}, opponent) &&
				!board.IsSquareUnderAttack(Square{2, rank}, opponent) {
				moves = append(moves, Square{2, rank})
			}
		}
	}
	return moves
}

func (c *CastlingPossibility) updateCastlingPossibility(lastMove *Move) {
	switch lastMove.ColoredPiece.Piece {
	case King:
		if lastMove.ColoredPiece.Color == White {
			c.IsWhiteKingsidePossible = false
			c.IsWhiteQueensidePossible = false
		} else {
			c.IsBlackKingsidePossible = false
			c.IsBlackQueensidePossible = false
		}
	case Rook:
		switch lastMove.OldSquare {
		case Square{File: 0, Rank: 0}:
			c.IsWhiteQueensidePossible = false
		case Square{File: 7, Rank: 0}:
			c.IsWhiteKingsidePossible = false
		case Square{File: 0, Rank: 7}:
			c.IsBlackQueensidePossible = false
		case Square{File: 7, Rank: 7}:
			c.IsBlackKingsidePossible = false
		}
	}
}

func (b *Board) castlingRookMove(movedPiece ColoredPiece, fromSquare Square, toSquare Square) {
	rank := fromSquare.Rank
	if toSquare.File == 2 {
		b.SetCell(0, rank, Empty, White)
		b.SetCell(3, rank, Rook, movedPiece.Color)
	} else {
		b.SetCell(7, rank, Empty, White)
		b.SetCell(5, rank, Rook, movedPiece.Color)
	}
}
