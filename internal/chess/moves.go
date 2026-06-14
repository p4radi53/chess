package chess

func (b *Board) pawnMoves(from Square, color Color) []Square {
	var moves []Square
	var direction int
	switch color {
	case White:
		direction = 1
	case Black:
		direction = -1
	}

	isFrontSquareEmpty := b.IsCellEmpty(from.File, from.Rank+direction)

	if isFrontSquareEmpty {
		moves = append(moves, Square{File: from.File, Rank: from.Rank + direction})
	}
	if (from.File-1 >= 0) && b.IsCellOccupiedByOpponent(from.File-1, from.Rank+direction, color) {
		moves = append(moves, Square{File: from.File - 1, Rank: from.Rank + direction})
	}
	if (from.File+1 < 8) && b.IsCellOccupiedByOpponent(from.File+1, from.Rank+direction, color) {
		moves = append(moves, Square{File: from.File + 1, Rank: from.Rank + direction})
	}

	// Initial double move
	if (color == White && from.Rank == 1) || (color == Black && from.Rank == 6) {
		if isFrontSquareEmpty && b.IsCellEmpty(from.File, from.Rank+2*direction) {
			moves = append(moves, Square{File: from.File, Rank: from.Rank + 2*direction})
		}
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

var bishopDirections = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var rookDirections = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var queenDirections = append(bishopDirections, rookDirections...)

func (b *Board) IsSquareAttackedByBishopQueenRook(square Square, byColor Color) bool {
	for _, dir := range queenDirections {
		for step := 1; step < 8; step++ {
			newFile := square.File + dir[0]*step
			newRank := square.Rank + dir[1]*step
			if !b.IsCellWithinBounds(newFile, newRank) {
				break
			}
			cell := b.GetCell(newFile, newRank)
			if cell.Piece != Empty {
				if cell.Color == byColor && (cell.Piece == Queen || (cell.Piece == Rook && (dir[0] == 0 || dir[1] == 0)) || (cell.Piece == Bishop && dir[0] != 0 && dir[1] != 0)) {
					return true
				}
				break
			}
		}
	}
	return false
}
func enemyColor(c Color) Color {
	if c == White {
		return Black
	}
	return White
}

func (b *Board) IsSquareUnderAttack(square Square, byColor Color) bool{
	if b.IsSquareAttackedByBishopQueenRook(square, byColor){
		return true
	}

	// horse
	for _, offset := range [][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}} {
		newFile := square.File + offset[0]
		newRank := square.Rank + offset[1]
		if b.IsCellWithinBounds(newFile, newRank){
			coloredPiece := b.GetCell(newFile, newRank)
			if (coloredPiece.Color == byColor && coloredPiece.Piece == Knight){
				return true
			}
		}
	}

	// pawn
	transform := 1
    if byColor == White{
        transform = -1
    }
    for _, fileDelta := range []int{-1, 1} {
        f, r := square.File+fileDelta, square.Rank+transform
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
		if !b.IsCellWithinBounds(newFile, newRank){
			continue;
		}
		if cp := b.GetCell(newFile, newRank); cp.Color == byColor && cp.Piece == King{
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
		if b.IsCellWithinBounds(newFile, newRank) && !b.IsCellOccupiedByOwnPiece(newFile, newRank, color) && !b.IsSquareUnderAttack(targetSquare, enemyColor(color)){
			moves = append(moves, Square{File: newFile, Rank: newRank})
		}
	}

	// Castle
	// wieze nie ruszone e
	// pola wolne
	// nie pod atackiem pole na ktorych rusza sie krol

	return moves
}

func (b *Board) LegalMoves(from Square) []Square {
	cell := b.GetCell(from.File, from.Rank)

	switch cell.Piece {
	case Empty:
		return nil
	case Pawn:
		return b.pawnMoves(from, cell.Color)
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
