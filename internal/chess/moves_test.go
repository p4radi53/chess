package chess

import (
	"slices"
	"testing"
)

func hasMove(moves []Square, file, rank int) bool {
	return slices.Contains(moves, Square{file, rank})
}

type moveTest struct {
	name     string
	setup    func(*Board)
	from     Square
	expected []Square
	excluded []Square
}

func runMoveTests(t *testing.T, tests []moveTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Board{}
			tt.setup(b)
			moves := b.LegalMoves(tt.from, Move{})
			for _, sq := range tt.expected {
				if !hasMove(moves, sq.File, sq.Rank) {
					t.Errorf("expected move to %v", sq)
				}
			}
			for _, sq := range tt.excluded {
				if hasMove(moves, sq.File, sq.Rank) {
					t.Errorf("should not have move to %v", sq)
				}
			}
		})
	}
}

func TestEnPassant(t *testing.T) {
	// White pawn on e5 (file=4, rank=4), black pawn just double-pushed from d7 to d5 (file=3)
	b := &Board{}
	b.SetCell(4, 4, Pawn, White)
	b.SetCell(3, 4, Pawn, Black)

	lastMove := Move{
		ColoredPiece: ColoredPiece{Piece: Pawn, Color: Black},
		OldSquare:    Square{3, 6},
		NewSquare:    Square{3, 4},
	}

	moves := b.LegalMoves(Square{4, 4}, lastMove)

	enPassantTarget := Square{3, 5}
	if !slices.Contains(moves, enPassantTarget) {
		t.Errorf("expected en passant capture to %v", enPassantTarget)
	}
}

func TestPawnMoves(t *testing.T) {
	tests := []moveTest{
		{
			name:     "single push",
			setup:    func(b *Board) { b.SetCell(4, 4, Pawn, White) },
			from:     Square{4, 4},
			expected: []Square{{4, 5}},
		},
		{
			name:     "double push from start",
			setup:    func(b *Board) { b.SetCell(4, 1, Pawn, White) },
			from:     Square{4, 1},
			expected: []Square{{4, 2}, {4, 3}},
		},
		{
			name: "double push on the edge",
			setup: func(b *Board) {
				b.SetCell(0, 1, Pawn, White)
			},
			from:     Square{0, 1},
			expected: []Square{{0, 2}, {0, 3}},
			excluded: []Square{},
		},
	}

	runMoveTests(t, tests)
}
