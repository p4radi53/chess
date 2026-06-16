package chess

type Move struct {
	ColoredPiece  ColoredPiece
	CapturedPiece ColoredPiece
	OldSquare     Square
	NewSquare     Square
}

var pawnDirection = []int{1, -1}       // White, Black
var pawnAttackDirection = []int{-1, 1} // Left, Right
var knightDirections = [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}}
var bishopDirections = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var rookDirections = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var queenDirections = append(bishopDirections, rookDirections...)

func (g *Game) pawnMoves(from Square, color Color, lastMove Move) []Square {
	var moves []Square

	if !g.Board.IsCellWithinBounds(from.File, from.Rank+pawnDirection[color]) {
		return moves
	}

	isFrontSquareEmpty := g.Board.IsCellEmpty(from.File, from.Rank+pawnDirection[color])

	if isFrontSquareEmpty {
		moves = append(moves, Square{File: from.File, Rank: from.Rank + pawnDirection[color]})
	}
	if (from.File-1 >= 0) && g.Board.IsCellOccupiedByOpponent(from.File-1, from.Rank+pawnDirection[color], color) {
		moves = append(moves, Square{File: from.File - 1, Rank: from.Rank + pawnDirection[color]})
	}
	if (from.File+1 < 8) && g.Board.IsCellOccupiedByOpponent(from.File+1, from.Rank+pawnDirection[color], color) {
		moves = append(moves, Square{File: from.File + 1, Rank: from.Rank + pawnDirection[color]})
	}

	// Initial double move
	if (color == White && from.Rank == 1) || (color == Black && from.Rank == 6) {
		if isFrontSquareEmpty && g.Board.IsCellEmpty(from.File, from.Rank+2*pawnDirection[color]) {
			moves = append(moves, Square{File: from.File, Rank: from.Rank + 2*pawnDirection[color]})
		}
	}

	// En passante
	if (lastMove.ColoredPiece.Piece == Pawn) &&
		(lastMove.NewSquare.Rank-lastMove.OldSquare.Rank == 2 ||
			lastMove.OldSquare.Rank-lastMove.NewSquare.Rank == 2) &&
		(lastMove.NewSquare.File == from.File-1 ||
			lastMove.NewSquare.File == from.File+1) &&
		lastMove.NewSquare.Rank == from.Rank {
		moves = append(moves, Square{File: lastMove.NewSquare.File,
			Rank: from.Rank + pawnDirection[color]})
	}

	return moves
}

func (g *Game) knightMoves(from Square, color Color) []Square {
	var moves []Square

	for _, offset := range knightDirections {
		newFile := from.File + offset[0]
		newRank := from.Rank + offset[1]
		if g.Board.IsCellWithinBounds(newFile, newRank) && !g.Board.IsCellOccupiedByOwnPiece(newFile, newRank, color) {
			moves = append(moves, Square{File: newFile, Rank: newRank})
		}
	}

	return moves
}

func (g *Game) slidingMoves(from Square, directions [][2]int, color Color) []Square {
	var moves []Square

	for _, dir := range directions {
		for step := 1; step < 8; step++ {
			newFile := from.File + dir[0]*step
			newRank := from.Rank + dir[1]*step
			if !g.Board.IsCellWithinBounds(newFile, newRank) {
				break
			}
			if g.Board.IsCellOccupiedByOwnPiece(newFile, newRank, color) {
				break
			}
			moves = append(moves, Square{File: newFile, Rank: newRank})
			if g.Board.IsCellOccupiedByOpponent(newFile, newRank, color) {
				break
			}
		}
	}
	return moves
}

func (g *Game) kingMoves(from Square, color Color) []Square {
	var moves []Square
	for _, offset := range queenDirections {
		newFile := from.File + offset[0]
		newRank := from.Rank + offset[1]
		targetSquare := Square{File: newFile, Rank: newRank}
		if g.Board.IsCellWithinBounds(newFile, newRank) && !g.Board.IsCellOccupiedByOwnPiece(newFile, newRank, color) && !g.Board.IsSquareUnderAttack(targetSquare, color.Opponent()) {
			moves = append(moves, Square{File: newFile, Rank: newRank})
		}
	}

	moves = append(moves, castlingMoves(g.CastlingPossibility, g.Board, color)...)

	return moves
}

func (g *Game) LegalMoves(from Square, lastMove Move) []Square {
	cell := g.Board.GetCell(from.File, from.Rank)

	switch cell.Piece {
	case Empty:
		return nil
	case Pawn:
		return g.pawnMoves(from, cell.Color, lastMove)
	case Knight:
		return g.knightMoves(from, cell.Color)
	case Bishop:
		return g.slidingMoves(from, bishopDirections, cell.Color)
	case Rook:
		return g.slidingMoves(from, rookDirections, cell.Color)
	case Queen:
		return g.slidingMoves(from, queenDirections, cell.Color)
	case King:
		return g.kingMoves(from, cell.Color)
	}
	return nil
}
