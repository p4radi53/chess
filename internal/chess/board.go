package chess

type Piece int

const (
	Empty Piece = iota
	King
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

type Color int

const (
	White Color = iota
	Black
)

type ColoredPiece struct {
	Piece Piece
	Color Color
}

type Board struct {
	Cells [8][8]ColoredPiece
}

type Square struct {
	File int `json:"file"`
	Rank int `json:"rank"`
}

func (b *Board) GetCell(file, rank int) ColoredPiece {
	return b.Cells[file][rank]
}

func (b *Board) IsCellEmpty(file, rank int) bool {
	return b.GetCell(file, rank).Piece == Empty
}

func (b *Board) IsCellOccupiedByOpponent(file, rank int, color Color) bool {
	cell := b.GetCell(file, rank)
	return cell.Piece != Empty && cell.Color != color
}

func (b *Board) IsCellOccupiedByOwnPiece(file, rank int, color Color) bool {
	cell := b.GetCell(file, rank)
	return cell.Piece != Empty && cell.Color == color
}

func (b *Board) IsCellWithinBounds(file, rank int) bool {
	return file >= 0 && file < 8 && rank >= 0 && rank < 8
}

func (b *Board) SetCell(file, rank int, piece Piece, color Color) {
	b.Cells[file][rank] = ColoredPiece{Piece: piece, Color: color}
}

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
	isLeftDiagonalSquareOccupiedByOpponent := b.IsCellOccupiedByOpponent(from.File-1, from.Rank+direction, color)
	isRightDiagonalSquareOccupiedByOpponent := b.IsCellOccupiedByOpponent(from.File+1, from.Rank+direction, color)

	if isFrontSquareEmpty {
		moves = append(moves, Square{File: from.File, Rank: from.Rank + direction})
	}
	if isLeftDiagonalSquareOccupiedByOpponent {
		moves = append(moves, Square{File: from.File - 1, Rank: from.Rank + direction})
	}
	if isRightDiagonalSquareOccupiedByOpponent {
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

func (b *Board) kingMoves(from Square, color Color) []Square {
	var moves []Square
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

func (b *Board) setupFirstPosition() {
	// Place pawns
	for i := range 8 {
		b.Cells[i][1] = ColoredPiece{Piece: Pawn, Color: White}
		b.Cells[i][6] = ColoredPiece{Piece: Pawn, Color: Black}
	}

	// Place other pieces
	pieces := []Piece{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, piece := range pieces {
		b.Cells[i][0] = ColoredPiece{Piece: piece, Color: White}
		b.Cells[i][7] = ColoredPiece{Piece: piece, Color: Black}
	}
}
