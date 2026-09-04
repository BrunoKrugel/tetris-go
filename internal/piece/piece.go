package piece

// Kind represents a tetromino type.
type Kind int

const KindEmpty Kind = -1

const (
	KindI Kind = 1 + iota
	KindO
	KindT
	KindS
	KindZ
	KindJ
	KindL
)

// Cell is an offset relative to the piece origin. Y grows downward.
type Cell struct {
	X, Y int
}

// Piece is a tetromino with a kind, rotation, and board position.
type Piece struct {
	Kind     Kind
	Rotation int // 0..3, normalized with %4
	X, Y     int // origin on the board
}

// baseCells holds the rotation-0 cells for each kind.
var baseCells = map[Kind][]Cell{
	KindI: {{0, 1}, {1, 1}, {2, 1}, {3, 1}},
	KindO: {{1, 0}, {2, 0}, {1, 1}, {2, 1}},
	KindT: {{0, 1}, {1, 0}, {1, 1}, {2, 1}},
	KindS: {{1, 0}, {2, 0}, {0, 1}, {1, 1}},
	KindZ: {{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	KindJ: {{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	KindL: {{2, 0}, {0, 1}, {1, 1}, {2, 1}},
}

// rotateCells applies 90° clockwise rotation n times to each cell.
func rotateCells(cells []Cell, n int) []Cell {
	r := make([]Cell, len(cells))
	copy(r, cells)
	for range n % 4 {
		for i, c := range r {
			r[i] = Cell{X: c.Y, Y: 3 - c.X}
		}
	}
	return r
}

// New creates a piece in spawn orientation at position (0,0).
func New(kind Kind) Piece {
	return Piece{Kind: kind}
}

// Cells returns the 4 cells for the current rotation, relative to the piece origin.
func (p Piece) Cells() []Cell {
	return rotateCells(baseCells[p.Kind], p.Rotation)
}

// Rotated returns a piece rotated 90° clockwise, normalized to 0..3.
func (p Piece) Rotated() Piece {
	return Piece{
		Kind:     p.Kind,
		Rotation: (p.Rotation + 1) % 4,
		X:        p.X,
		Y:        p.Y,
	}
}

// AllKinds returns all playable tetromino kinds.
func AllKinds() []Kind {
	return []Kind{KindI, KindO, KindT, KindS, KindZ, KindJ, KindL}
}

func (k Kind) String() string {
	switch k {
	case KindI:
		return "I"
	case KindO:
		return "O"
	case KindT:
		return "T"
	case KindS:
		return "S"
	case KindZ:
		return "Z"
	case KindJ:
		return "J"
	case KindL:
		return "L"
	default:
		return "?"
	}
}
