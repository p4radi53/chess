import { useState, useEffect } from "react";
import type { GameState, Square } from "./types";

const API = "/api";

const PIECE_LETTER = ["", "K", "Q", "R", "B", "N", "P"];

function pieceImage(piece: number, color: number): string | null {
  if (piece === 0) return null;
  const c = color === 0 ? "w" : "b";
  return `/pieces/${c}${PIECE_LETTER[piece]}.svg`;
}

const ranks = [7, 6, 5, 4, 3, 2, 1, 0];
const files = [0, 1, 2, 3, 4, 5, 6, 7];

export default function Board() {
  const [game, setGame] = useState<GameState | null>(null);
  const [selected, setSelected] = useState<Square | null>(null);
  const [legalMoves, setLegalMoves] = useState<Square[]>([]);

  useEffect(() => {
    fetch(`${API}/game`)
      .then((r) => r.json())
      .then(setGame);
  }, []);

  async function handleSquareClick(square: Square) {
    if (!game) return;
    const cell = game.Board.Cells[square.file][square.rank];

    if (selected === null) {
      if (cell.Piece === 0) return;
      if (cell.Color !== game.CurrentTurn) return;
      setSelected(square);
      const res = await fetch(
        `${API}/legal-moves?file=${square.file}&rank=${square.rank}`,
      );
      if (res.ok) setLegalMoves(await res.json());
    } else {
      const isLegal = legalMoves.some(
        (m) => m.file === square.file && m.rank === square.rank,
      );
      if (!isLegal) {
        setSelected(null);
        setLegalMoves([]);
        return;
      }
      setSelected(null);
      setLegalMoves([]);
      const res = await fetch(`${API}/move`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          from_file: selected.file,
          from_rank: selected.rank,
          to_file: square.file,
          to_rank: square.rank,
        }),
      });
      if (res.ok) {
        setGame(await res.json());
      } else {
        console.error("Move failed:", await res.json());
      }
    }
  }

  async function handleReset() {
    const res = await fetch(`${API}/reset`, { method: "POST" });
    if (res.ok) setGame(await res.json());
  }

  if (!game) {
    return (
      <div className="flex h-screen items-center justify-center bg-zinc-900">
        <span className="text-zinc-400">Loading...</span>
      </div>
    );
  }

  // cell size: fills viewport on mobile, capped at 64px on desktop
  const cellSize = "min(11vw, 64px)";
  const pieceSize = "min(8.5vw, 52px)";

  return (
    <div className="flex h-screen items-center justify-center bg-zinc-900">
      <div className="flex flex-col items-center gap-2">
        <table
          className="border-2 border-zinc-600"
          style={{ borderCollapse: "collapse" }}
        >
          <tbody>
            {ranks.map((rank) => (
              <tr key={rank}>
                <td
                  className="text-right text-zinc-400 select-none"
                  style={{
                    paddingRight: "clamp(2px, 1vw, 8px)",
                    fontSize: cellSize,
                    width: cellSize,
                  }}
                >
                  {rank + 1}
                </td>
                {files.map((file) => {
                  const cell = game.Board.Cells[file][rank];
                  const isLight = (file + rank) % 2 !== 0;
                  const isSelected =
                    selected?.file === file && selected?.rank === rank;
                  const isLegal = legalMoves.some(
                    (m) => m.file === file && m.rank === rank,
                  );
                  const img = pieceImage(cell.Piece, cell.Color);

                  return (
                    <td
                      key={file}
                      onClick={() => handleSquareClick({ file, rank })}
                      style={{
                        width: cellSize,
                        height: cellSize,
                        background: isSelected
                          ? "#4fc3f7"
                          : isLegal
                            ? "#a8d8a8"
                            : isLight
                              ? "#c8bea0"
                              : "#82644b",
                        textAlign: "center",
                        verticalAlign: "middle",
                        cursor: "pointer",
                        userSelect: "none",
                        outline: isSelected ? "3px solid #0ea5e9" : undefined,
                        outlineOffset: isSelected ? "-3px" : undefined,
                      }}
                    >
                      <div
                        style={{
                          width: cellSize,
                          height: cellSize,
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                        }}
                      >
                        {img && (
                          <img
                            src={img}
                            style={{ width: pieceSize, height: pieceSize }}
                            alt=""
                            draggable={false}
                          />
                        )}
                      </div>
                    </td>
                  );
                })}
              </tr>
            ))}
          </tbody>
        </table>
        <div className="text-sm text-zinc-300">
          {game.CurrentTurn === 0 ? "White" : "Black"} to move
          {game.IsWhiteInCheck && " · White is in check"}
          {game.IsBlackInCheck && " · Black is in check"}
        </div>
        <button
          onClick={handleReset}
          className="mt-2 px-4 py-1 text-sm text-zinc-300 bg-zinc-700 hover:bg-zinc-600 rounded"
        >
          Reset
        </button>
        <div
          className="flex text-sm text-zinc-400 select-none"
          style={{ paddingLeft: "clamp(12px, 3vw, 28px)" }}
        >
          {["a", "b", "c", "d", "e", "f", "g", "h"].map((l) => (
            <div key={l} style={{ width: cellSize, textAlign: "center" }}>
              {l}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
