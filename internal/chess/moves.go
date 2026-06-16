package chess

type Move struct {
	ColoredPiece ColoredPiece
	OldSquare    Square
	NewSquare    Square
}

var colorDirection = []int{1, -1} // White, Black
var knightDirections = [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}} 
var bishopDirections = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var rookDirections = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var queenDirections = append(bishopDirections, rookDirections...)

func (b *Board) pawnMoves(from Square, color Color, lastMove Move) []Square {
	var moves []Square

	if !b.IsCellWithinBounds(newFile, newRank + colorDirection[color]){
		return moves
	} 

	isFrontSquareEmpty := b.IsCellEmpty(from.File, from.Rank+colorDirection[color])

	if isFrontSquareEmpty {
		moves = append(moves, Square{File: from.File, Rank: from.Rank + colorDirection[color]})
	}
	if (from.File-1 >= 0) && b.IsCellOccupiedByOpponent(from.File-1, from.Rank + colorDirection[color], color) {
		moves = append(moves, Square{File: from.File - 1, Rank: from.Rank + colorDirection[color]})
	}
	if (from.File+1 < 8) && b.IsCellOccupiedByOpponent(from.File+1, from.Rank + colorDirection[color], color) {
		moves = append(moves, Square{File: from.File + 1, Rank: from.Rank + colorDirection[color]})
	}

	// Initial double move
	if (color == White && from.Rank == 1) || (color == Black && from.Rank == 6) {
		if isFrontSquareEmpty && b.IsCellEmpty(from.File, from.Rank + 2 * colorDirection[color]) {
			moves = append(moves, Square{File: from.File, Rank: from.Rank + 2*colorDirection[color]})
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
			Rank: from.Rank + direction})
	}

	return moves
}

func (b *Board) knightMoves(from Square, color Color) []Square {
	var moves []Square

	for _, offset := range [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}} {
		newFile := from.File + offset[0]
		newRank := from.Rank + offset[1]
		if b.IsCellWithinBounds(newFile, newRank) && !b.IsCellOccupiedByOwnPiece(newFile, newRank, color) {
			moves = append(moves, Square{File: newFile, Rank: newRank})
		}
	}

	return moves
}

func (b *Board) slidingMoves(from Square, directions [][2]int, color Color) []Square {
	var moves []Square

	for _, dir := range directions {
		for step := 1; step < 8; step++ {
			newFile := from.File + dir[0]*step
			newRank := from.Rank + dir[1]*step
			if !b.IsCellWithinBounds(newFile, newRank) {
				break
			}
			if b.IsCellOccupiedByOwnPiece(newFile, newRank, color) {
				break
			}
			moves = append(moves, Square{File: newFile, Rank: newRank})
			if b.IsCellOccupiedByOpponent(newFile, newRank, color) {
				break
			}
		}
	}
	return moves
}



func (b *Board) IsSquareAttackedByBishopQueenRook(square Square, attackingColor Color) bool {
	for _, dir := range queenDirections {
		for step := 1; step < 8; step++ {
			newFile := square.File + dir[0]*step
			newRank := square.Rank + dir[1]*step
			if !b.IsCellWithinBounds(newFile, newRank) {
				break
			}
			cell := b.GetCell(newFile, newRank)
			if cell.Piece != Empty {
				if cell.Color == attackingColor && (cell.Piece == Queen || (cell.Piece == Rook && (dir[0] == 0 || dir[1] == 0)) || (cell.Piece == Bishop && dir[0] != 0 && dir[1] != 0)) {
					return true
				}
				break
			}
		}
	}
	return false
}

func (b *Board) IsSquareAttackedByKnight(square Square, attackingColor Color) bool {
for _, offset := knightDirections{
		newFile := square.File + offset[0]
		newRank := square.Rank + offset[1]
		if b.IsCellWithinBounds(newFile, newRank) {
			coloredPiece := b.GetCell(newFile, newRank)
			if coloredPiece.Color == attackingColor && coloredPiece.Piece == Knight {
				return true
			}
		}
}
return false
}
func (b *Board) IsSquareUnderAttack(square Square, attackingColor Color) bool {
	if b.IsSquareAttackedByBishopQueenRook(square, attackingColor)
	// horse
	for _, offset := range [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}} {
		newFile := square.File + offset[0]
		newRank := square.Rank + offset[1]
		if b.IsCellWithinBounds(newFile, newRank) {
			coloredPiece := b.GetCell(newFile, newRank)
			if coloredPiece.Color == attackingColor && coloredPiece.Piece == Knight {
				return true
			}
		coloredPiece := b.GetCell(newFile, newRank)
		if b.IsCellWithinBounds(newFile, newRank) && (coloredPiece.Color == byColor && coloredPiece.Piece == Knight){
			return true
		}
	}

	// pawn
	for _, fileDelta := range []int{-1, 1} {
		f, r := square.File+fileDelta, square.Rank + colorDirection[attackingColor]
		if b.IsCellWithinBounds(f, r) {
			if cp := b.GetCell(f, r); cp.Piece == Pawn && cp.Color == byColor {
				return true
			}
		}
	}

	// king
	for _, offset := range [][2]int{{1, 1}, {1, 0}, {1, -1}, {0, 1}, {0, -1}, {-1, 1}, {-1, 0}, {-1, -1}} {
		newFile := square.File + offset[0]
		newRank := square.Rank + offset[1]
		if !b.IsCellWithinBounds(newFile, newRank) {
			continue
		}
		if cp := b.GetCell(newFile, newRank); b.IsCellWithinBounds(newFile, newRank) && cp.Color == byColor && cp.Piece == King{
			return true
		}
	}
	return false
}

func (b *Board) kingMoves(from Square, color Color) []Square {
	var moves []Square
	for _, offset := range [][2]int{{1, 1}, {1, 0}, {1, -1}, {0, 1}, {0, -1}, {-1, 1}, {-1, 0}, {-1, -1}} {
		newFile := from.File + offset[0]
		newRank := from.Rank + offset[1]
		targetSquare := Square{File: newFile, Rank: newRank}
		if b.IsCellWithinBounds(newFile, newRank) && !b.IsCellOccupiedByOwnPiece(newFile, newRank, color) && !b.IsSquareUnderAttack(targetSquare, color.Opponent()) {
			moves = append(moves, Square{File: newFile, Rank: newRank})
		}
	}

	// TODO: add castling
	// TODO: discovered checks

	return moves
}

func (b *Board) LegalMoves(from Square, lastMove Move) []Square {
	cell := b.GetCell(from.File, from.Rank)

	switch cell.Piece {
	case Empty:
		return nil
	case Pawn:
		return b.pawnMoves(from, cell.Color, lastMove)
	case Knight:
		return b.knightMoves(from, cell.Color)
	case Bishop:
		return b.slidingMoves(from, bishopDirections, cell.Color)
	case Rook:
		return b.slidingMoves(from, rookDirections, cell.Color)
	case Queen:
		return b.slidingMoves(from, queenDirections, cell.Color)
	case King:
		return b.kingMoves(from, cell.Color)
	}
	return nil
}
