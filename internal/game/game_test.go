package game

import (
	"testing"
	"tetris-go/internal/board"
	"tetris-go/internal/piece"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestGame(kinds ...piece.Kind) *Game {
	return New(NewDeterministic(kinds...))
}

func Test_New(t *testing.T) {
	t.Run("Should start at level 1 with score 0 and current plus next pieces", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindO)

		assert.Equal(t, 1, g.Level)
		assert.Equal(t, 0, g.Score)
		assert.Equal(t, 0, g.Lines)
		assert.Equal(t, StatePlaying, g.State)
		assert.Equal(t, piece.KindI, g.Current.Kind)
		assert.Equal(t, piece.KindO, g.Next.Kind)
	})
}

func Test_Move(t *testing.T) {
	t.Run("Should decrease X when moving left", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindI)
		before := g.Current.X
		g.MoveLeft()
		assert.Equal(t, before-1, g.Current.X)
	})

	t.Run("Should not move outside the left wall", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.Current.X = 0
		for range 20 {
			g.MoveLeft()
		}
		for _, c := range g.Current.Cells() {
			assert.GreaterOrEqual(t, g.Current.X+c.X, 0)
		}
	})

	t.Run("Should not move outside the right wall", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		for range 20 {
			g.MoveRight()
		}
		maxX := board.Width - 1
		for _, c := range g.Current.Cells() {
			assert.LessOrEqual(t, g.Current.X+c.X, maxX)
		}
	})

	t.Run("Should move down and increase score by 1 on soft drop", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.Current.Y = 0
		g.SoftDrop()
		assert.Equal(t, 1, g.Current.Y)
		assert.Equal(t, 1, g.Score)
	})
}

func Test_Ghost(t *testing.T) {
	t.Run("Should return the lowest valid position on an empty board", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.Current.X = 0
		g.Current.Y = 0
		ghost := g.Ghost()
		next := ghost
		next.Y++
		assert.False(t, g.Board.CanPlace(next))
	})

	t.Run("Should stop immediately above existing blocks", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.Board.Set(piece.Piece{Kind: piece.KindO, X: 0, Y: 18})
		g.Current.X = 0
		g.Current.Y = 0
		ghost := g.Ghost()
		for _, c := range ghost.Cells() {
			assert.NotEqual(t, piece.KindO, g.Board.Grid[ghost.Y+c.Y][c.X])
		}
	})
}

func Test_HardDrop(t *testing.T) {
	t.Run("Should lock the piece and spawn the next one", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindI)
		g.HardDrop()
		assert.Equal(t, piece.KindI, g.Current.Kind)
		assert.NotEqual(t, 0, g.Score)
	})

	t.Run("Should stop above existing blocks", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.Board.Set(piece.Piece{Kind: piece.KindO, X: 4, Y: 14})
		g.Current.X = 4
		g.Current.Y = 0
		g.HardDrop()
		occupied := 0
		for y := range board.Height {
			for x := range board.Width {
				if g.Board.Grid[y][x] != piece.KindEmpty {
					occupied++
				}
			}
		}
		assert.Equal(t, 8, occupied)
	})
}

func Test_Rotate(t *testing.T) {
	t.Run("Should succeed in open space", func(t *testing.T) {
		g := newTestGame(piece.KindT, piece.KindT)
		g.Current.X = 4
		g.Current.Y = 8
		g.Rotate()
		assert.Equal(t, 1, g.Current.Rotation)
	})

	t.Run("Should apply a wall kick near the left wall", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindI)
		g.Current.X = -2
		g.Current.Y = 8
		g.Rotate()
		assert.True(t, g.Board.CanPlace(g.Current))
	})

	t.Run("Should be rejected when no kick position fits", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		for _, y := range []int{9, 10, 11} {
			for x := range board.Width {
				g.Board.Grid[y][x] = piece.KindJ
			}
		}
		g.Current.X = 4
		g.Current.Y = 8
		g.Rotate()
		assert.Equal(t, 0, g.Current.Rotation)
	})
}

func Test_Hold(t *testing.T) {
	t.Run("Should swap current with hold and spawn next", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindO, piece.KindT)
		g.Hold()
		assert.Equal(t, piece.KindI, g.Held.Kind)
		assert.Equal(t, piece.KindO, g.Current.Kind)
		assert.Equal(t, piece.KindT, g.Next.Kind)
	})

	t.Run("Should reject a second hold before locking", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindO, piece.KindT)
		g.Hold()
		held := g.Held.Kind
		g.Hold()
		assert.Equal(t, held, g.Held.Kind)
	})

	t.Run("Should swap with the held piece when hold is occupied", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindO, piece.KindT)
		g.Hold()
		g.HardDrop()
		g.Hold()
		assert.Equal(t, piece.KindT, g.Held.Kind)
		assert.Equal(t, piece.KindI, g.Current.Kind)
	})

	t.Run("Should enable hold again after a lock", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindO, piece.KindT)
		g.Hold()
		require.False(t, g.CanHold)
		g.HardDrop()
		assert.True(t, g.CanHold)
	})
}

func Test_Scoring(t *testing.T) {
	tests := []struct {
		name     string
		lines    int
		level    int
		expected int
	}{
		{"Should score 100 for one line at level 1", 1, 1, 100},
		{"Should score 300 for two lines at level 1", 2, 1, 300},
		{"Should score 500 for three lines at level 1", 3, 1, 500},
		{"Should score 800 for a Tetris at level 1", 4, 1, 800},
		{"Should score 200 for one line at level 2", 1, 2, 200},
		{"Should score 1600 for a Tetris at level 2", 4, 2, 1600},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGame(piece.KindI, piece.KindI)
			g.Level = tt.level
			g.addScore(tt.lines)
			assert.Equal(t, tt.expected, g.Score)
		})
	}
}

func Test_LevelProgression(t *testing.T) {
	t.Run("Should stay at level 1 below 10 lines", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindI)
		g.Lines = 9
		g.Level = 1 + g.Lines/LinesPerLevel
		assert.Equal(t, 1, g.Level)
	})

	t.Run("Should reach level 2 at 10 lines", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindI)
		g.Lines = 10
		g.Level = 1 + g.Lines/LinesPerLevel
		assert.Equal(t, 2, g.Level)
	})

	t.Run("Should increase fall speed per level", func(t *testing.T) {
		assert.Greater(t, FallSpeedForLevel(1), FallSpeedForLevel(2))
		assert.Greater(t, FallSpeedForLevel(2), FallSpeedForLevel(5))
	})

	t.Run("Should never go below the minimum speed", func(t *testing.T) {
		assert.Equal(t, 0.1, FallSpeedForLevel(100))
	})
}

func Test_GameOver(t *testing.T) {
	t.Run("Should end the game when spawn is blocked", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		for y := range board.Height {
			for x := range board.Width {
				g.Board.Grid[y][x] = piece.KindJ
			}
		}
		g.Current = g.newPiece()
		g.HardDrop()
		assert.Equal(t, StateGameOver, g.State)
	})

	t.Run("Should ignore updates after game over", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.State = StateGameOver
		score := g.Score
		g.Update(10)
		g.MoveLeft()
		assert.Equal(t, score, g.Score)
	})
}

func Test_Pause(t *testing.T) {
	t.Run("Should pause and stop falling", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.TogglePause()
		assert.Equal(t, StatePaused, g.State)
		y := g.Current.Y
		g.Update(10)
		assert.Equal(t, y, g.Current.Y)
	})

	t.Run("Should disable movement while paused", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.TogglePause()
		x := g.Current.X
		g.MoveLeft()
		assert.Equal(t, x, g.Current.X)
	})

	t.Run("Should resume playing", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindO)
		g.TogglePause()
		g.TogglePause()
		assert.Equal(t, StatePlaying, g.State)
	})
}

func Test_Restart(t *testing.T) {
	t.Run("Should reset the game completely", func(t *testing.T) {
		g := newTestGame(piece.KindI, piece.KindI)
		g.HardDrop()
		g.HardDrop()
		g.TogglePause()
		g.Restart()
		assert.Equal(t, 0, g.Score)
		assert.Equal(t, 1, g.Level)
		assert.Equal(t, 0, g.Lines)
		assert.Equal(t, StatePlaying, g.State)
		empty := true
		for y := range board.Height {
			for x := range board.Width {
				if g.Board.Grid[y][x] != piece.KindEmpty {
					empty = false
				}
			}
		}
		assert.True(t, empty)
	})
}

func Test_Update(t *testing.T) {
	t.Run("Should apply gravity and lock on floor contact", func(t *testing.T) {
		g := newTestGame(piece.KindO, piece.KindI)
		g.Update(g.FallSpeed)
		assert.Equal(t, SpawnY+1, g.Current.Y)
		for range 100 {
			g.Update(g.FallSpeed)
		}
		assert.NotEqual(t, piece.KindO, g.Current.Kind)
	})
}

func Test_SevenBag(t *testing.T) {
	t.Run("Should yield each kind exactly once per seven draws", func(t *testing.T) {
		bag := NewSevenBag(42)
		counts := map[piece.Kind]int{}
		for range 7 {
			counts[bag.Next()]++
		}
		for _, kind := range piece.AllKinds() {
			assert.Equal(t, 1, counts[kind], "kind %s", kind)
		}
	})
}
