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

func (c Color) Opponent() Color {
	return c ^ 1
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
