package chess

import (
	"slices"
	"testing"
)

func gameWithClearBoard() *Game {
	return &Game{
		CastlingPossibility: CastlingPossibility{
			IsWhiteKingsidePossible:  true,
			IsWhiteQueensidePossible: true,
			IsBlackKingsidePossible:  true,
			IsBlackQueensidePossible: true,
		},
	}
}

func assertRookMove(t *testing.T, b *Board, fromFile, fromRank, toFile, toRank int, color Color) {
	t.Helper()
	if b.GetCell(fromFile, fromRank).Piece != Empty {
		t.Errorf("expected {%d,%d} to be empty after castling", fromFile, fromRank)
	}
	if cell := b.GetCell(toFile, toRank); cell.Piece != Rook || cell.Color != color {
		t.Errorf("expected %v rook on {%d,%d}, got %v", color, toFile, toRank, cell)
	}
}

// ── castlingMoves ────────────────────────────────────────────────────────────

func TestCastlingMovesAllowedSides(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		kingFile int
		rank     int
		rookFile int
		wantFile int
	}{
		{"white kingside", White, 4, 0, 7, 6},
		{"white queenside", White, 4, 0, 0, 2},
		{"black kingside", Black, 4, 7, 7, 6},
		{"black queenside", Black, 4, 7, 0, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gameWithClearBoard()
			g.Board.SetCell(tt.kingFile, tt.rank, King, tt.color)
			g.Board.SetCell(tt.rookFile, tt.rank, Rook, tt.color)

			moves := castlingMoves(g.CastlingPossibility, g.Board, tt.color)

			if !slices.Contains(moves, Square{tt.wantFile, tt.rank}) {
				t.Errorf("expected castling to {%d,%d}, got %v", tt.wantFile, tt.rank, moves)
			}
		})
	}
}

func TestCastlingMovesBlockedByPiece(t *testing.T) {
	g := gameWithClearBoard()
	g.Board.SetCell(4, 0, King, White)
	g.Board.SetCell(7, 0, Rook, White)
	g.Board.SetCell(6, 0, Knight, White)

	moves := castlingMoves(g.CastlingPossibility, g.Board, White)

	if slices.Contains(moves, Square{6, 0}) {
		t.Errorf("kingside castling should be blocked by piece on f1")
	}
}

func TestCastlingMovesBlockedByAttack(t *testing.T) {
	g := gameWithClearBoard()
	g.Board.SetCell(4, 0, King, White)
	g.Board.SetCell(7, 0, Rook, White)
	g.Board.SetCell(5, 7, Rook, Black) // attacks f1

	moves := castlingMoves(g.CastlingPossibility, g.Board, White)

	if slices.Contains(moves, Square{6, 0}) {
		t.Errorf("kingside castling should be blocked when f1 is under attack")
	}
}

func TestCastlingMovesNoRightsNoMoves(t *testing.T) {
	g := gameWithClearBoard()
	g.CastlingPossibility = CastlingPossibility{}
	g.Board.SetCell(4, 0, King, White)
	g.Board.SetCell(7, 0, Rook, White)
	g.Board.SetCell(0, 0, Rook, White)

	moves := castlingMoves(g.CastlingPossibility, g.Board, White)

	if len(moves) != 0 {
		t.Errorf("expected no castling moves when rights revoked, got %v", moves)
	}
}

// ── castlingRookMove ─────────────────────────────────────────────────────────

func TestCastlingRookMove(t *testing.T) {
	tests := []struct {
		name                     string
		color                    Color
		rookFile, rank           int
		kingFrom, kingTo         Square
		wantRookFrom, wantRookTo int
	}{
		{"white kingside", White, 7, 0, Square{4, 0}, Square{6, 0}, 7, 5},
		{"white queenside", White, 0, 0, Square{4, 0}, Square{2, 0}, 0, 3},
		{"black kingside", Black, 7, 7, Square{4, 7}, Square{6, 7}, 7, 5},
		{"black queenside", Black, 0, 7, Square{4, 7}, Square{2, 7}, 0, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gameWithClearBoard()
			g.Board.SetCell(tt.rookFile, tt.rank, Rook, tt.color)
			g.Board.castlingRookMove(ColoredPiece{King, tt.color}, tt.kingFrom, tt.kingTo)
			assertRookMove(t, &g.Board, tt.wantRookFrom, tt.rank, tt.wantRookTo, tt.rank, tt.color)
		})
	}
}

// ── updateCastlingPossibility ────────────────────────────────────────────────

func TestUpdateCastlingPossibilityKingMove(t *testing.T) {
	c := CastlingPossibility{true, true, true, true}
	c.updateCastlingPossibility(&Move{ColoredPiece: ColoredPiece{King, White}})

	if c.IsWhiteKingsidePossible || c.IsWhiteQueensidePossible {
		t.Errorf("white castling rights should be revoked after king moves")
	}
	if !c.IsBlackKingsidePossible || !c.IsBlackQueensidePossible {
		t.Errorf("black castling rights should be unaffected")
	}
}

func TestUpdateCastlingPossibilityRookMove(t *testing.T) {
	tests := []struct {
		name      string
		rookFrom  Square
		wantFalse *bool
		wantTrue  *bool
	}{
		{"white kingside rook", Square{7, 0}, nil, nil},
		{"white queenside rook", Square{0, 0}, nil, nil},
		{"black kingside rook", Square{7, 7}, nil, nil},
		{"black queenside rook", Square{0, 7}, nil, nil},
	}

	// use a function instead of field pointers — cleaner
	type check struct {
		name     string
		rookFrom Square
		assertFn func(CastlingPossibility) (bool, bool) // (shouldBeFalse, shouldBeTrue)
	}
	checks := []check{
		{"white kingside rook", Square{7, 0}, func(c CastlingPossibility) (bool, bool) {
			return c.IsWhiteKingsidePossible, c.IsWhiteQueensidePossible
		}},
		{"white queenside rook", Square{0, 0}, func(c CastlingPossibility) (bool, bool) {
			return c.IsWhiteQueensidePossible, c.IsWhiteKingsidePossible
		}},
		{"black kingside rook", Square{7, 7}, func(c CastlingPossibility) (bool, bool) {
			return c.IsBlackKingsidePossible, c.IsBlackQueensidePossible
		}},
		{"black queenside rook", Square{0, 7}, func(c CastlingPossibility) (bool, bool) {
			return c.IsBlackQueensidePossible, c.IsBlackKingsidePossible
		}},
	}
	_ = tests
	for _, tt := range checks {
		t.Run(tt.name, func(t *testing.T) {
			c := CastlingPossibility{true, true, true, true}
			c.updateCastlingPossibility(&Move{
				ColoredPiece: ColoredPiece{Rook, White},
				OldSquare:    tt.rookFrom,
			})
			revoked, kept := tt.assertFn(c)
			if revoked {
				t.Errorf("expected right to be revoked")
			}
			if !kept {
				t.Errorf("expected other right to be unaffected")
			}
		})
	}
}

func TestUpdateCastlingPossibilityUnrelatedMove(t *testing.T) {
	c := CastlingPossibility{true, true, true, true}
	c.updateCastlingPossibility(&Move{ColoredPiece: ColoredPiece{Pawn, White}})

	if !c.IsWhiteKingsidePossible || !c.IsWhiteQueensidePossible ||
		!c.IsBlackKingsidePossible || !c.IsBlackQueensidePossible {
		t.Errorf("castling rights should be unchanged after pawn move")
	}
}

// ── MakeMove integration ─────────────────────────────────────────────────────

func TestMakeMoveWhiteKingsideCastling(t *testing.T) {
	g := gameWithClearBoard()
	g.CurrentTurn = White
	g.Board.SetCell(4, 0, King, White)
	g.Board.SetCell(7, 0, Rook, White)
	g.WhiteKingPosition = Square{4, 0}
	g.BlackKingPosition = Square{4, 7}

	if err := g.MakeMove(4, 0, 6, 0); err != nil {
		t.Fatalf("expected kingside castling to succeed: %v", err)
	}
	if cell := g.Board.GetCell(6, 0); cell.Piece != King || cell.Color != White {
		t.Errorf("expected white king on g1, got %v", cell)
	}
	assertRookMove(t, &g.Board, 7, 0, 5, 0, White)
}

func TestMakeMoveWhiteQueensideCastling(t *testing.T) {
	g := gameWithClearBoard()
	g.CurrentTurn = White
	g.Board.SetCell(4, 0, King, White)
	g.Board.SetCell(0, 0, Rook, White)
	g.WhiteKingPosition = Square{4, 0}
	g.BlackKingPosition = Square{4, 7}

	if err := g.MakeMove(4, 0, 2, 0); err != nil {
		t.Fatalf("expected queenside castling to succeed: %v", err)
	}
	if cell := g.Board.GetCell(2, 0); cell.Piece != King || cell.Color != White {
		t.Errorf("expected white king on c1, got %v", cell)
	}
	assertRookMove(t, &g.Board, 0, 0, 3, 0, White)
}
