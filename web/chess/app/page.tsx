import Board from "./Board";

const API = "http://localhost:8080";

export type Cell = {
  File: number;
  Rank: number;
  Piece: number;
  Color: number;
};

export type BoardState = {
  Cells: Cell[][];
};

export default async function Home() {
  const res = await fetch(`${API}/board`, { cache: "no-store" });
  const board: BoardState = await res.json();

  return <Board initialBoard={board} />;
}
