# Chess

A chess game with a Go backend and a Next.js frontend.

## Stack

- **Backend** — Go 1.26, Gin HTTP server
- **Frontend** — Next.js 16, React, Tailwind CSS
- **Pieces** — cburnett SVG set (lichess)

## Project structure

```
cmd/server/        # entrypoint
internal/chess/    # game logic (board, moves, rules)
internal/server/   # HTTP handlers
web/chess/         # Next.js app
```

## Running

**Backend** (from project root):
```bash
go run cmd/server/main.go
```

**Frontend** (from `web/chess`):
```bash
bun dev
# or
npm run dev
```

Then open `http://localhost:3000`.
