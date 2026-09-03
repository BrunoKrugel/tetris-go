package board

import "tetris-go/internal/piece"

const (
	Width  = 10
	Height = 20
)

type Board struct {
	Grid [Height][Width]piece.Kind
}

func New() *Board {
	b := &Board{}
	for y := range Height {
		for x := range Width {
			b.Grid[y][x] = piece.KindEmpty
		}
	}
	return b
}

func (b *Board) IsInside(x, y int) bool {
	return x >= 0 && x < Width && y >= 0 && y < Height
}

func (b *Board) CanPlace(p piece.Piece) bool {
	for _, c := range p.Cells() {
		bx, by := p.X+c.X, p.Y+c.Y
		if bx < 0 || bx >= Width || by >= Height {
			return false
		}
		if by < 0 {
			continue
		}
		if b.Grid[by][bx] != piece.KindEmpty {
			return false
		}
	}
	return true
}

func (b *Board) Set(p piece.Piece) {
	for _, c := range p.Cells() {
		bx, by := p.X+c.X, p.Y+c.Y
		if by >= 0 {
			b.Grid[by][bx] = p.Kind
		}
	}
}

func (b *Board) ClearLines() int {
	var cleared int
	y := 0
	for y < Height {
		full := true
		for x := range Width {
			if b.Grid[y][x] == piece.KindEmpty {
				full = false
				break
			}
		}
		if full {
			for row := y; row > 0; row-- {
				b.Grid[row] = b.Grid[row-1]
			}
			for x := range Width {
				b.Grid[0][x] = piece.KindEmpty
			}
			cleared++
		} else {
			y++
		}
	}
	return cleared
}
