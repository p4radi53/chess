package chess

import (
	"fmt"
	"slices"
)

type Game struct {
	Board               Board
	CurrentTurn         Color
	WhiteKingPosition   Square
	BlackKingPosition   Square
	IsWhiteInCheck      bool
	IsBlackInCheck      bool
	Moves               []Move
	CastlingPossibility CastlingPossibility
}

func NewGame() *Game {
	game := &Game{
		CurrentTurn: White,
		CastlingPossibility: CastlingPossibility{
			IsWhiteKingsidePossible:  true,
			IsWhiteQueensidePossible: true,
			IsBlackKingsidePossible:  true,
			IsBlackQueensidePossible: true,
		},
	}
	game.Board.setupFirstPosition()
	game.WhiteKingPosition = Square{4, 0}
	game.BlackKingPosition = Square{4, 7}
	return game
}

func (g *Game) MakeMove(fromFile, fromRank, toFile, toRank int) error {
	sourceCell := g.Board.GetCell(fromFile, fromRank)
	movedColoredPiece := ColoredPiece{Piece: sourceCell.Piece, Color: sourceCell.Color}

	if sourceCell.Piece == Empty {
		return fmt.Errorf("no piece at source square")
	}
	if sourceCell.Color != g.CurrentTurn {
		return fmt.Errorf("not %v's turn", sourceCell.Color)
	}

	if !slices.Contains(g.Position().LegalMoves(Square{fromFile, fromRank}), Square{toFile,
		toRank}) {
		return fmt.Errorf("not a legal move for the piece")
	}

	targetCell := g.Board.GetCell(toFile, toRank)
	var removedPiece ColoredPiece
	if targetCell.Piece != Empty && targetCell.Color == g.CurrentTurn {
		return fmt.Errorf("cannot capture own piece")
	} else if targetCell.Piece != Empty {
		removedPiece = ColoredPiece{Piece: targetCell.Piece, Color: movedColoredPiece.Color.Opponent()}
	}
	g.Board.SetCell(fromFile, fromRank, Empty, White)
	g.Board.SetCell(toFile, toRank, sourceCell.Piece, sourceCell.Color)

	if sourceCell.Piece == King && (toFile-fromFile > 1 || fromFile-toFile > 1) {
		g.Board.castlingRookMove(movedColoredPiece, Square{fromFile, fromRank}, Square{toFile, toRank})
	}

	// Remove the captured pawn for en passant
	if sourceCell.Piece == Pawn && toFile != fromFile && targetCell.Piece == Empty {
		removedPiece = ColoredPiece{Piece: Pawn, Color: movedColoredPiece.Color.Opponent()}
		g.Board.SetCell(toFile, fromRank, Empty, White)
	}

	currentMove := Move{ColoredPiece: movedColoredPiece, CapturedPiece: removedPiece, OldSquare: Square{File: fromFile, Rank: fromRank}, NewSquare: Square{File: toFile, Rank: toRank}}
	g.Moves = append(g.Moves, currentMove)

	g.updateCheckFlags(currentMove)
	if movedColoredPiece.Piece == King {
		g.updateKingPosition(currentMove)
	}
	g.CastlingPossibility.updateCastlingPossibility(&currentMove)
	g.CurrentTurn = (g.CurrentTurn + 1) % 2
	return nil
}

func (g *Game) LastMove() Move {
	if len(g.Moves) == 0 {
		return Move{}
	}
	return g.Moves[len(g.Moves)-1]
}

func (g *Game) Position() Position {
	return Position{
		Board:       g.Board,
		Castling:    g.CastlingPossibility,
		SideToMove:  g.CurrentTurn,
		LastMove:    g.LastMove(),
		KingSquares: [2]Square{g.WhiteKingPosition, g.BlackKingPosition},
	}
}

func (g *Game) updateKingPosition(move Move) {
	movedColor := move.ColoredPiece.Color
	if movedColor == White {
		g.WhiteKingPosition = move.NewSquare
	} else if movedColor == Black {
		g.BlackKingPosition = move.NewSquare
	}
}

func (g *Game) updateCheckFlags(move Move) {
	var kingPos Square
	if move.ColoredPiece.Color == White {
		kingPos = g.BlackKingPosition
	} else {
		kingPos = g.WhiteKingPosition
	}

	isCheck := g.Board.IsSquareUnderAttack(kingPos, move.ColoredPiece.Color.Opponent())
	if move.ColoredPiece.Color == White {
		g.IsBlackInCheck = isCheck
	} else {
		g.IsWhiteInCheck = isCheck
	}
}
