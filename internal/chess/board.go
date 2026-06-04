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
	File int
	Rank int
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

func (b *Board) SetCell(file, rank int, piece Piece, color Color) {
	b.Cells[file][rank] = ColoredPiece{Piece: piece, Color: color}
}

func (b *Board) isInBounds(file, rank int) bool {
	return file >= 0 && file < 8 && rank >= 0 && rank < 8
}

func (b *Board) pawnMoves(from Square, color Color) []Square {
	var moves []Square
	var direction int
	if color == White {
		direction = 1
	} else if color == Black {
		direction = -1
	}

	// no sense in calculating moves if the pawn is already on the last rank
	if !b.isInBounds(from.File, from.Rank+direction) {
		return moves
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
	return moves
}

func (b *Board) slidingMoves(from Square, directions [][2]int, color Color) []Square {
	var moves []Square
	return moves
}

func (b *Board) kingMoves(from Square, color Color) []Square {
	var moves []Square
	return moves
}

var bishopDirections = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var rookDirections = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var queenDirections = append(bishopDirections, rookDirections...)

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
