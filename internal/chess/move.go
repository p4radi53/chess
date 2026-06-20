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
