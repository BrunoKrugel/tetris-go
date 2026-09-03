package board

import (
	"testing"
	"tetris-go/internal/piece"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoard(t *testing.T) {
	t.Run("TC-BOARD-001 New board has all cells empty", func(t *testing.T) {
		b := New()

		for y := range Height {
			for x := range Width {
				assert.Equal(t, piece.KindEmpty, b.Grid[y][x], "cell (%d,%d)", x, y)
			}
		}
	})

	t.Run("TC-BOARD-002 IsInside returns true for valid positions", func(t *testing.T) {
		b := New()

		assert.True(t, b.IsInside(0, 0))
		assert.True(t, b.IsInside(9, 0))
		assert.True(t, b.IsInside(0, 19))
		assert.True(t, b.IsInside(9, 19))
		assert.True(t, b.IsInside(5, 10))
	})

	t.Run("TC-BOARD-003 IsInside returns false for out-of-bounds positions", func(t *testing.T) {
		b := New()

		assert.False(t, b.IsInside(-1, 0))
		assert.False(t, b.IsInside(10, 0))
		assert.False(t, b.IsInside(0, 20))
		assert.False(t, b.IsInside(10, 20))
	})

	t.Run("TC-BOARD-004 IsInside returns false for above-board y<0", func(t *testing.T) {
		b := New()

		assert.False(t, b.IsInside(0, -1))
		assert.False(t, b.IsInside(5, -5))
	})
}

func TestCollision(t *testing.T) {
	t.Run("TC-COLLISION-001 returns false when piece touches floor", func(t *testing.T) {
		b := New()
		p := piece.New(piece.KindI)
		p.Y = Height

		assert.False(t, b.CanPlace(p))
	})

	t.Run("TC-COLLISION-002 returns false when piece overlaps existing block", func(t *testing.T) {
		b := New()
		p := piece.New(piece.KindI)
		p.X, p.Y = 0, 18
		require.True(t, b.CanPlace(p))
		b.Set(p)

		assert.False(t, b.CanPlace(p))
	})

	t.Run("TC-COLLISION-003 returns false when piece collides with side wall", func(t *testing.T) {
		b := New()
		p := piece.New(piece.KindI)
		p.X = -1

		assert.False(t, b.CanPlace(p))
	})
}

func TestClearLines(t *testing.T) {
	t.Run("TC-LINE-001 returns 0 when no lines are full", func(t *testing.T) {
		b := New()

		assert.Equal(t, 0, b.ClearLines())
	})

	t.Run("TC-LINE-002 clears 1 line", func(t *testing.T) {
		b := New()
		fillRow(b, 19)

		assert.Equal(t, 1, b.ClearLines())
	})

	t.Run("TC-LINE-003 clears 2 lines", func(t *testing.T) {
		b := New()
		fillRow(b, 18)
		fillRow(b, 19)

		assert.Equal(t, 2, b.ClearLines())
	})

	t.Run("TC-LINE-004 clears 3 lines", func(t *testing.T) {
		b := New()
		fillRow(b, 17)
		fillRow(b, 18)
		fillRow(b, 19)

		assert.Equal(t, 3, b.ClearLines())
	})

	t.Run("TC-LINE-005 clears 4 lines", func(t *testing.T) {
		b := New()
		fillRow(b, 16)
		fillRow(b, 17)
		fillRow(b, 18)
		fillRow(b, 19)

		assert.Equal(t, 4, b.ClearLines())
	})

	t.Run("TC-LINE-006 blocks above full line fall down correctly", func(t *testing.T) {
		b := New()
		// Place a block at (0, 18)
		b.Grid[18][0] = piece.KindI
		// Fill row 19 completely
		fillRow(b, 19)

		cleared := b.ClearLines()

		assert.Equal(t, 1, cleared)
		assert.Equal(t, piece.KindI, b.Grid[19][0], "block from row 18 should fall to row 19")
		assert.Equal(t, piece.KindEmpty, b.Grid[18][0], "row 18 should be empty after shift")
	})

	t.Run("TC-LINE-007 nearly full board clears nothing", func(t *testing.T) {
		b := New()
		// Fill every cell except one per row, so no row is complete
		for y := range Height {
			for x := range Width {
				b.Grid[y][x] = piece.KindI
			}
			b.Grid[y][0] = piece.KindEmpty // leave one gap per row
		}

		assert.Equal(t, 0, b.ClearLines())
	})
}

func fillRow(b *Board, y int) {
	for x := range Width {
		b.Grid[y][x] = piece.KindI
	}
}
