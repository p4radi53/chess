package chess

import (
	"fmt"
	"slices"
)

type Game struct {
	Board         Board
	CurrentTurn   Color
	MoveCount     int
	RemovedPieces []Piece
}

func NewGame() *Game {
	game := &Game{
		CurrentTurn: White,
		MoveCount:   0,
	}
	game.Board.setupFirstPosition()
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
	g.MoveCount++
	g.CurrentTurn = (g.CurrentTurn + 1) % 2
	return nil
}
