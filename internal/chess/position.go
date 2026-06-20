package chess

type Position struct {
	Board       Board
	Castling    CastlingPossibility
	SideToMove  Color
	LastMove    Move
	KingSquares [2]Square
}

func (p Position) pawnMoves(from Square, color Color) []Square {
	var moves []Square

	if !p.Board.IsCellWithinBounds(from.File, from.Rank+pawnDirection[color]) {
		return moves
	}

	isFrontSquareEmpty := p.Board.IsCellEmpty(from.File, from.Rank+pawnDirection[color])

	if isFrontSquareEmpty {
		moves = append(moves, Square{File: from.File, Rank: from.Rank + pawnDirection[color]})
	}
	if (from.File-1 >= 0) && p.Board.IsCellOccupiedByOpponent(from.File-1, from.Rank+pawnDirection[color], color) {
		moves = append(moves, Square{File: from.File - 1, Rank: from.Rank + pawnDirection[color]})
	}
	if (from.File+1 < 8) && p.Board.IsCellOccupiedByOpponent(from.File+1, from.Rank+pawnDirection[color], color) {
		moves = append(moves, Square{File: from.File + 1, Rank: from.Rank + pawnDirection[color]})
	}

	if (color == White && from.Rank == 1) || (color == Black && from.Rank == 6) {
		if isFrontSquareEmpty && p.Board.IsCellEmpty(from.File, from.Rank+2*pawnDirection[color]) {
			moves = append(moves, Square{File: from.File, Rank: from.Rank + 2*pawnDirection[color]})
		}
	}

	if p.LastMove.ColoredPiece.Piece == Pawn &&
		Abs(p.LastMove.NewSquare.Rank-p.LastMove.OldSquare.Rank) == 2 &&
		Abs(p.LastMove.NewSquare.File-from.File) == 1 &&
		p.LastMove.NewSquare.Rank == from.Rank {
		moves = append(moves, Square{File: p.LastMove.NewSquare.File,
			Rank: from.Rank + pawnDirection[color]})
	}

	return moves
}

func (p Position) knightMoves(from Square, color Color) []Square {
	var moves []Square

	for _, offset := range knightDirections {
		newFile := from.File + offset[0]
		newRank := from.Rank + offset[1]
		if p.Board.IsCellWithinBounds(newFile, newRank) && !p.Board.IsCellOccupiedByOwnPiece(newFile, newRank, color) {
			moves = append(moves, Square{File: newFile, Rank: newRank})
		}
	}

	return moves
}

func (p Position) slidingMoves(from Square, directions [][2]int, color Color) []Square {
	var moves []Square

	for _, dir := range directions {
		for step := 1; step < 8; step++ {
			newFile := from.File + dir[0]*step
			newRank := from.Rank + dir[1]*step
			if !p.Board.IsCellWithinBounds(newFile, newRank) {
				break
			}
			if p.Board.IsCellOccupiedByOwnPiece(newFile, newRank, color) {
				break
			}
			moves = append(moves, Square{File: newFile, Rank: newRank})
			if p.Board.IsCellOccupiedByOpponent(newFile, newRank, color) {
				break
			}
		}
	}
	return moves
}

func (p Position) kingMoves(from Square, color Color) []Square {
	var moves []Square
	for _, offset := range queenDirections {
		newFile := from.File + offset[0]
		newRank := from.Rank + offset[1]
		targetSquare := Square{File: newFile, Rank: newRank}
		if p.Board.IsCellWithinBounds(newFile, newRank) && !p.Board.IsCellOccupiedByOwnPiece(newFile, newRank, color) && !p.Board.IsSquareUnderAttack(targetSquare, color.Opponent()) {
			moves = append(moves, Square{File: newFile, Rank: newRank})
		}
	}

	moves = append(moves, castlingMoves(p.Castling, p.Board, color)...)

	return moves
}

func (p Position) LegalMoves(from Square) []Square {
	cell := p.Board.GetCell(from.File, from.Rank)

	switch cell.Piece {
	case Empty:
		return nil
	case Pawn:
		return p.pawnMoves(from, cell.Color)
	case Knight:
		return p.knightMoves(from, cell.Color)
	case Bishop:
		return p.slidingMoves(from, bishopDirections, cell.Color)
	case Rook:
		return p.slidingMoves(from, rookDirections, cell.Color)
	case Queen:
		return p.slidingMoves(from, queenDirections, cell.Color)
	case King:
		return p.kingMoves(from, cell.Color)
	}
	return nil
}
