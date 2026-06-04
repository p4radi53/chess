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
			moves := b.LegalMoves(tt.from)
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
	}

	runMoveTests(t, tests)
}
