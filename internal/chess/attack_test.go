package chess

import "testing"

type attackTest struct {
	name           string
	setup          func(*Board)
	square         Square
	attackingColor Color
	expected       bool
}

func runAttackTests(t *testing.T, tests []attackTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Board{}
			tt.setup(b)
			got := b.IsSquareUnderAttack(tt.square, tt.attackingColor)
			if got != tt.expected {
				t.Errorf("IsSquareUnderAttack(%v, %v) = %v, want %v", tt.square, tt.attackingColor, got, tt.expected)
			}
		})
	}
}

// ── Rook ────────────────────────────────────────────────────────────────────

func TestRookAttacks(t *testing.T) {
	runAttackTests(t, []attackTest{
		{
			name: "rook attacks along rank",
			setup: func(b *Board) {
				b.SetCell(0, 4, Rook, Black)
			},
			square:         Square{4, 4},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "rook attacks along file",
			setup: func(b *Board) {
				b.SetCell(3, 0, Rook, Black)
			},
			square:         Square{3, 7},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "rook blocked by own piece",
			setup: func(b *Board) {
				b.SetCell(0, 4, Rook, Black)
				b.SetCell(2, 4, Pawn, Black)
			},
			square:         Square{4, 4},
			attackingColor: Black,
			expected:       false,
		},
		{
			name: "rook blocked by opponent piece",
			setup: func(b *Board) {
				b.SetCell(0, 4, Rook, Black)
				b.SetCell(2, 4, Pawn, White)
			},
			square:         Square{4, 4},
			attackingColor: Black,
			expected:       false,
		},
		{
			name: "rook does not attack diagonally",
			setup: func(b *Board) {
				b.SetCell(0, 0, Rook, Black)
			},
			square:         Square{3, 3},
			attackingColor: Black,
			expected:       false,
		},
	})
}

// ── Bishop ───────────────────────────────────────────────────────────────────

func TestBishopAttacks(t *testing.T) {
	runAttackTests(t, []attackTest{
		{
			name: "bishop attacks diagonally",
			setup: func(b *Board) {
				b.SetCell(0, 0, Bishop, Black)
			},
			square:         Square{5, 5},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "bishop attacks anti-diagonal",
			setup: func(b *Board) {
				b.SetCell(7, 0, Bishop, Black)
			},
			square:         Square{3, 4},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "bishop blocked on diagonal",
			setup: func(b *Board) {
				b.SetCell(0, 0, Bishop, Black)
				b.SetCell(2, 2, Pawn, White)
			},
			square:         Square{5, 5},
			attackingColor: Black,
			expected:       false,
		},
		{
			name: "bishop does not attack along rank",
			setup: func(b *Board) {
				b.SetCell(0, 4, Bishop, Black)
			},
			square:         Square{4, 4},
			attackingColor: Black,
			expected:       false,
		},
	})
}

// ── Queen ────────────────────────────────────────────────────────────────────

func TestQueenAttacks(t *testing.T) {
	runAttackTests(t, []attackTest{
		{
			name: "queen attacks along rank",
			setup: func(b *Board) {
				b.SetCell(0, 3, Queen, White)
			},
			square:         Square{7, 3},
			attackingColor: White,
			expected:       true,
		},
		{
			name: "queen attacks diagonally",
			setup: func(b *Board) {
				b.SetCell(1, 1, Queen, White)
			},
			square:         Square{4, 4},
			attackingColor: White,
			expected:       true,
		},
		{
			name: "queen blocked on rank",
			setup: func(b *Board) {
				b.SetCell(0, 3, Queen, White)
				b.SetCell(3, 3, Pawn, White)
			},
			square:         Square{7, 3},
			attackingColor: White,
			expected:       false,
		},
	})
}

// ── Knight ───────────────────────────────────────────────────────────────────

func TestKnightAttacks(t *testing.T) {
	runAttackTests(t, []attackTest{
		{
			name: "knight attacks L-shape (+2+1)",
			setup: func(b *Board) {
				b.SetCell(3, 3, Knight, Black)
			},
			square:         Square{5, 4},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "knight attacks L-shape (-1+2)",
			setup: func(b *Board) {
				b.SetCell(3, 3, Knight, Black)
			},
			square:         Square{2, 5},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "knight does not attack adjacent square",
			setup: func(b *Board) {
				b.SetCell(3, 3, Knight, Black)
			},
			square:         Square{4, 3},
			attackingColor: Black,
			expected:       false,
		},
		{
			name: "knight does not attack diagonal",
			setup: func(b *Board) {
				b.SetCell(3, 3, Knight, Black)
			},
			square:         Square{4, 4},
			attackingColor: Black,
			expected:       false,
		},
		{
			name: "knight jumps over pieces",
			setup: func(b *Board) {
				b.SetCell(3, 3, Knight, Black)
				b.SetCell(3, 4, Pawn, White)
				b.SetCell(4, 3, Pawn, White)
			},
			square:         Square{5, 4},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "knight near board edge does not wrap",
			setup: func(b *Board) {
				b.SetCell(0, 0, Knight, Black)
			},
			square:         Square{7, 1}, // would need wrap to reach
			attackingColor: Black,
			expected:       false,
		},
	})
}

// ── Pawn ─────────────────────────────────────────────────────────────────────

func TestPawnAttacks(t *testing.T) {
	runAttackTests(t, []attackTest{
		{
			name: "white pawn attacks left-diagonal",
			setup: func(b *Board) {
				b.SetCell(3, 3, Pawn, White)
			},
			square:         Square{2, 4},
			attackingColor: White,
			expected:       true,
		},
		{
			name: "white pawn attacks right-diagonal",
			setup: func(b *Board) {
				b.SetCell(3, 3, Pawn, White)
			},
			square:         Square{4, 4},
			attackingColor: White,
			expected:       true,
		},
		{
			name: "white pawn does not attack forward",
			setup: func(b *Board) {
				b.SetCell(3, 3, Pawn, White)
			},
			square:         Square{3, 4},
			attackingColor: White,
			expected:       false,
		},
		{
			name: "white pawn does not attack backward",
			setup: func(b *Board) {
				b.SetCell(3, 3, Pawn, White)
			},
			square:         Square{2, 2},
			attackingColor: White,
			expected:       false,
		},
		{
			name: "black pawn attacks downward-diagonal",
			setup: func(b *Board) {
				b.SetCell(3, 4, Pawn, Black)
			},
			square:         Square{2, 3},
			attackingColor: Black,
			expected:       true,
		},
		{
			name: "black pawn does not attack upward",
			setup: func(b *Board) {
				b.SetCell(3, 4, Pawn, Black)
			},
			square:         Square{2, 5},
			attackingColor: Black,
			expected:       false,
		},
	})
}

// ── King ─────────────────────────────────────────────────────────────────────

func TestKingAttacks(t *testing.T) {
	runAttackTests(t, []attackTest{
		{
			name: "king attacks adjacent square horizontally",
			setup: func(b *Board) {
				b.SetCell(3, 3, King, White)
			},
			square:         Square{4, 3},
			attackingColor: White,
			expected:       true,
		},
		{
			name: "king attacks adjacent square diagonally",
			setup: func(b *Board) {
				b.SetCell(3, 3, King, White)
			},
			square:         Square{4, 4},
			attackingColor: White,
			expected:       true,
		},
		{
			name: "king does not attack two squares away",
			setup: func(b *Board) {
				b.SetCell(3, 3, King, White)
			},
			square:         Square{5, 3},
			attackingColor: White,
			expected:       false,
		},
		{
			name: "king does not attack two squares diagonally",
			setup: func(b *Board) {
				b.SetCell(3, 3, King, White)
			},
			square:         Square{5, 5},
			attackingColor: White,
			expected:       false,
		},
	})
}

// ── Mixed / edge cases ───────────────────────────────────────────────────────

func TestNoAttacker(t *testing.T) {
	runAttackTests(t, []attackTest{
		{
			name:           "empty board",
			setup:          func(b *Board) {},
			square:         Square{4, 4},
			attackingColor: Black,
			expected:       false,
		},
		{
			name: "only friendly pieces on board",
			setup: func(b *Board) {
				b.SetCell(0, 4, Rook, White)
				b.SetCell(4, 0, Bishop, White)
			},
			square:         Square{4, 4},
			attackingColor: Black,
			expected:       false,
		},
	})
}
