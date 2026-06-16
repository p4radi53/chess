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
  IsWhiteInCheck: boolean;
  IsBlackInCheck: boolean;
};

export type Square = {
  file: number;
  rank: number;
};
