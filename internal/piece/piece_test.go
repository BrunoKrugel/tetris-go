package piece

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllKinds_HaveFourCells(t *testing.T) {
	for _, kind := range AllKinds() {
		t.Run("Should have 4 cells for "+kind.String(), func(t *testing.T) {
			cells := New(kind).Cells()
			require.Len(t, cells, 4)
		})
	}
}

func TestAllKinds_CellsWithinBoundingBox(t *testing.T) {
	for _, kind := range AllKinds() {
		for rot := range 4 {
			t.Run(kind.String()+" rot="+string(rune('0'+rot)), func(t *testing.T) {
				p := Piece{Kind: kind, Rotation: rot}
				for _, c := range p.Cells() {
					assert.True(t, c.X >= 0 && c.X <= 3, "X=%d out of [0,3]", c.X)
					assert.True(t, c.Y >= 0 && c.Y <= 3, "Y=%d out of [0,3]", c.Y)
				}
			})
		}
	}
}

func TestPiece_FourRotationsReturnToOriginal(t *testing.T) {
	for _, kind := range AllKinds() {
		t.Run("Should return to original after 4 rotations for "+kind.String(), func(t *testing.T) {
			p := New(kind)
			original := p.Cells()
			for range 4 {
				p = p.Rotated()
			}
			assert.Equal(t, original, p.Cells())
		})
	}
}

func TestPiece_SpawnOrientations(t *testing.T) {
	tests := []struct {
		expected []Cell
		kind     Kind
	}{
		{kind: KindI, expected: []Cell{{0, 1}, {1, 1}, {2, 1}, {3, 1}}},
		{kind: KindO, expected: []Cell{{1, 0}, {2, 0}, {1, 1}, {2, 1}}},
		{kind: KindT, expected: []Cell{{0, 1}, {1, 0}, {1, 1}, {2, 1}}},
		{kind: KindS, expected: []Cell{{1, 0}, {2, 0}, {0, 1}, {1, 1}}},
		{kind: KindZ, expected: []Cell{{0, 0}, {1, 0}, {1, 1}, {2, 1}}},
		{kind: KindJ, expected: []Cell{{0, 0}, {0, 1}, {1, 1}, {2, 1}}},
		{kind: KindL, expected: []Cell{{2, 0}, {0, 1}, {1, 1}, {2, 1}}},
	}

	for _, tt := range tests {
		t.Run("Should have correct spawn shape for "+tt.kind.String(), func(t *testing.T) {
			p := New(tt.kind)
			assert.ElementsMatch(t, tt.expected, p.Cells())
		})
	}
}
