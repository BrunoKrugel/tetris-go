# Tetris (Go + Ebiten)

A simple, standard-rules Tetris designed to be easy to read and play — big
pieces, big text, no clutter. Builds to a standalone Windows `.exe`.

## Run

```bash
make run        # or: go run .
```

## Controls

| Key   | Action                          |
| ----- | ------------------------------- |
| ← →   | Move left / right               |
| ↓     | Soft drop                       |
| ↑     | Rotate                          |
| Space | Hard drop                       |
| C     | Hold piece                      |
| P     | Pause / resume                  |
| R     | Restart (after game over, too)  |

## Rules

- 10×20 board, 7 tetrominoes, 7-bag randomizer, wall kicks on rotation.
- Score: 1 line = 100 × level, 2 = 300 ×, 3 = 500 ×, 4 (Tetris) = 800 ×.
- Soft drop +1/cell, hard drop +2/cell. Level up every 10 lines (faster fall).

## Development

```bash
make test       # unit tests
make race       # race detector
make cover      # coverage report
make windows    # cross-compile tetris.exe (Windows x64, com ícone e sem console)
```

Architecture: pure game engine in `internal/` (no Ebiten dependency) —
`internal/piece` (tetromino shapes/rotation), `internal/board` (grid,
collision, line clearing), `internal/game` (state machine, scoring, levels,
hold/ghost/hard drop) — and `main.go` (Ebiten rendering + keyboard input).
Tests use a deterministic piece generator, no randomness.

See `plan.md` for the full test plan and acceptance criteria.
