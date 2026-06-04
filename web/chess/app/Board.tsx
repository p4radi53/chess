"use client";

import { useState } from "react";
import Image from "next/image";
import { BoardState } from "./page";

const API = "http://localhost:8080";

const PIECE_LETTER = ["", "K", "Q", "R", "B", "N", "P"];

function pieceImage(piece: number, color: number): string | null {
  if (piece === 0) return null;
  const c = color === 0 ? "w" : "b";
  return `/pieces/${c}${PIECE_LETTER[piece]}.svg`;
}

const ranks = [7, 6, 5, 4, 3, 2, 1, 0];
const files = [0, 1, 2, 3, 4, 5, 6, 7];

interface Square {
  File: number;
  Rank: number;
}

interface Coordinates {
  file: number;
  rank: number;
}

export default function Board({ initialBoard }: { initialBoard: BoardState }) {
  const [board, setBoard] = useState<BoardState>(initialBoard);
  const [selected, setSelected] = useState<{ coordinates: Coordinates } | null>(
    null,
  );
  const [legalMoves, setLegalMoves] = useState<Square[]>([]);

  async function handleSquareClick(coordinates: Coordinates) {
    const cell = board.Cells[coordinates.file][coordinates.rank];

    if (selected === null) {
      if (cell.Piece === 0) return;
      setSelected({ coordinates });
      const res = await fetch(
        `${API}/legal-moves?file=${coordinates.file}&rank=${coordinates.rank}`,
      );
      if (res.ok) setLegalMoves(await res.json());
    } else {
      const from = selected;
      setSelected(null);
      setLegalMoves([]);
      const res = await fetch(`${API}/move`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          from_file: from.coordinates.file,
          from_rank: from.coordinates.rank,
          to_file: coordinates.file,
          to_rank: coordinates.rank,
        }),
      });
      if (res.ok) {
        setBoard(await res.json());
      } else {
        const err = await res.json();
        console.error("Move failed:", err);
      }
    }
  }

  async function handleReset() {
    const res = await fetch(`${API}/reset`, { method: "POST" });
    if (res.ok) setBoard(await res.json());
  }

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
                <td className="pr-2 text-right text-sm text-zinc-400 select-none w-5">
                  {rank + 1}
                </td>
                {files.map((file) => {
                  const cell = board.Cells[file][rank];
                  const isLight = (file + rank) % 2 !== 0;
                  const isSelected =
                    selected?.coordinates.file === file &&
                    selected?.coordinates.rank === rank;
                  const isLegal = legalMoves.some(
                    (m: Square) => m.File === file && m.Rank === rank,
                  );
                  const img = pieceImage(cell.Piece, cell.Color);

                  return (
                    <td
                      key={file}
                      onClick={() => handleSquareClick({ file, rank })}
                      style={{
                        width: 64,
                        height: 64,
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
                          width: 64,
                          height: 64,
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                        }}
                      >
                        {img && (
                          <Image
                            src={img}
                            width={52}
                            height={52}
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
        <button
          onClick={handleReset}
          className="mt-2 px-4 py-1 text-sm text-zinc-300 bg-zinc-700 hover:bg-zinc-600 rounded"
        >
          Reset
        </button>
        <div
          className="flex text-sm text-zinc-400 select-none"
          style={{ paddingLeft: 28 }}
        >
          {["a", "b", "c", "d", "e", "f", "g", "h"].map((l) => (
            <div key={l} style={{ width: 64, textAlign: "center" }}>
              {l}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
