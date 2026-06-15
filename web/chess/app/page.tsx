import Board from "./Board";

const API = process.env.API_URL!;

export type ColoredPiece = {
  Piece: number;
  Color: number;
};

export type BoardState = {
  Cells: ColoredPiece[][];
};

export type GameState = {
  Board: BoardState;
  CurrentTurn: number;
  MoveCount: number;
  IsWhiteInCheck: boolean;
  IsBlackInCheck: boolean;
};

export default async function Home() {
  const res = await fetch(`${API}/game`, { cache: "no-store" });
  const game: GameState = await res.json();

  return <Board initialGame={game} />;
}
