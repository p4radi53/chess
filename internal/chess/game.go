package chess

import (
	"fmt"
	"slices"
)

type Game struct {
	Board             Board
	CurrentTurn       Color
	MoveCount         int
	RemovedPieces     []Piece
	WhiteKingPosition Square
	BlackKingPosition Square
	IsWhiteInCheck    bool
	IsBlackInCheck    bool
}

func NewGame() *Game {
	game := &Game{
		CurrentTurn: White,
		MoveCount:   0,
	}
	game.Board.setupFirstPosition()
	game.WhiteKingPosition = Square{4, 0}
	game.BlackKingPosition = Square{4, 7}
	return game
}

func (g *Game) MakeMove(fromFile, fromRank, toFile, toRank int) error {
	sourceCell := g.Board.GetCell(fromFile, fromRank)
	if sourceCell.Piece == Empty {
		return fmt.Errorf("no piece at source square")
	}
	if sourceCell.Color != g.CurrentTurn {
		return fmt.Errorf("not %v's turn", sourceCell.Color)
	}

	if !slices.Contains(g.Board.LegalMoves(Square{fromFile, fromRank}), Square{toFile,
		toRank}) {
		return fmt.Errorf("not a legal move for the piece")
	}

	targetCell := g.Board.GetCell(toFile, toRank)
	if targetCell.Piece != Empty && targetCell.Color == g.CurrentTurn {
		return fmt.Errorf("cannot capture own piece")
	} else if targetCell.Piece != Empty {
		g.RemovedPieces = append(g.RemovedPieces, targetCell.Piece)
	}
	g.Board.SetCell(fromFile, fromRank, Empty, White)
	g.Board.SetCell(toFile, toRank, sourceCell.Piece, sourceCell.Color)
	if g.detectCheck(toFile, toRank) {
		switch g.CurrentTurn {
		case White:
			g.IsBlackInCheck = true
		case Black:
			g.IsWhiteInCheck = true
		}
	}
	g.MoveCount++
	g.CurrentTurn = (g.CurrentTurn + 1) % 2
	return nil
}

func (g *Game) detectCheck(toFile, toRank int) bool {
	targetCell := g.Board.GetCell(toFile, toRank)

	movedColor := targetCell.Color
	opponentColor := (movedColor + 1) % 2

	var kingPos Square
	if movedColor == White {
		kingPos = g.BlackKingPosition
	} else {
		kingPos = g.WhiteKingPosition
	}
	switch targetCell.Piece {
	case Pawn, Knight:
		return slices.Contains(g.Board.LegalMoves(Square{toFile, toRank}), kingPos)
	}

	return g.Board.IsSquareAttackedByBishopQueenRook(kingPos, opponentColor)
}
