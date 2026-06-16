package chess

func (b *Board) isSquareAttackedByBishopQueenRook(square Square, attackingColor Color) bool {
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

func (b *Board) isSquareAttackedByKnight(square Square, attackingColor Color) bool {
	for _, offset := range knightDirections {
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

func (b *Board) isSquareAttackedByPawn(square Square, attackingcolor Color) bool {
	for _, offset := range pawnDirection {
		newFile, newRank := square.File+offset, square.Rank+colorDirection[attackingcolor]
		if b.IsCellWithinBounds(newFile, newRank) {
			if coloredPiece := b.GetCell(newFile, newRank); coloredPiece.Piece == Pawn && coloredPiece.Color == attackingcolor {
				return true
			}
		}
	}
	return false
}

func (b *Board) isSquareAttackedByKing(square Square, attackingcolor Color) bool {
	for _, offset := range queenDirections {
		newFile := square.File + offset[0]
		newRank := square.Rank + offset[1]
		if !b.IsCellWithinBounds(newFile, newRank) {
			continue
		}
		if coloredPiece := b.GetCell(newFile, newRank); coloredPiece.Piece == King && coloredPiece.Color == attackingcolor {
			return true
		}
	}
	return false
}

func (b *Board) IsSquareUnderAttack(square Square, attackingcolor Color) bool {
	if b.isSquareAttackedByBishopQueenRook(square, attackingcolor) {
		return true
	}
	if b.isSquareAttackedByKnight(square, attackingcolor) {
		return true
	}
	if b.isSquareAttackedByPawn(square, attackingcolor) {
		return true
	}
	if b.isSquareAttackedByKing(square, attackingcolor) {
		return true
	}
	return false
}
