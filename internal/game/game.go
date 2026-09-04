package game

import (
	"math/rand/v2"
	"tetris-go/internal/board"
	"tetris-go/internal/piece"
)

type State int

const (
	StatePlaying State = iota
	StatePaused
	StateGameOver
)

const (
	LinesPerLevel = 10
	SpawnX        = 3
	SpawnY        = -2
)

type PieceGenerator interface {
	Next() piece.Kind
}

type SevenBag struct {
	rnd   *rand.Rand
	queue []piece.Kind
}

func NewSevenBag(seed uint64) *SevenBag {
	//nolint:gosec // PCG is fine for shuffling tetrominoes, no security relevance
	return &SevenBag{rnd: rand.New(rand.NewPCG(seed, seed))}
}

func (b *SevenBag) refill() {
	kinds := piece.AllKinds()
	b.rnd.Shuffle(len(kinds), func(i, j int) {
		kinds[i], kinds[j] = kinds[j], kinds[i]
	})
	b.queue = kinds
}

func (b *SevenBag) Next() piece.Kind {
	if len(b.queue) == 0 {
		b.refill()
	}
	k := b.queue[0]
	b.queue = b.queue[1:]
	return k
}

type DeterministicGenerator struct {
	kinds []piece.Kind
	index int
}

func NewDeterministic(kinds ...piece.Kind) *DeterministicGenerator {
	return &DeterministicGenerator{kinds: kinds}
}

func (g *DeterministicGenerator) Next() piece.Kind {
	k := g.kinds[g.index%len(g.kinds)]
	g.index++
	return k
}

type Game struct {
	generator PieceGenerator
	Board     *board.Board
	Current   piece.Piece
	Next      piece.Piece
	Held      piece.Piece
	Score     int
	Lines     int
	Level     int
	State     State
	FallSpeed float64
	fallTimer float64
	CanHold   bool
}

func New(gen PieceGenerator) *Game {
	g := &Game{
		Board:     board.New(),
		Held:      piece.Piece{Kind: piece.KindEmpty},
		Level:     1,
		FallSpeed: 1.0,
		State:     StatePlaying,
		generator: gen,
	}
	g.Next = g.newPiece()
	g.spawn()
	return g
}

func (g *Game) newPiece() piece.Piece {
	p := piece.New(g.generator.Next())
	p.X = SpawnX
	p.Y = SpawnY
	return p
}

func (g *Game) spawn() {
	g.Current = g.Next
	g.Next = g.newPiece()
	g.CanHold = true
	if !g.Board.CanPlace(g.Current) {
		g.State = StateGameOver
	}
}

func (g *Game) tryMove(dx, dy int) bool {
	if g.State != StatePlaying {
		return false
	}
	candidate := g.Current
	candidate.X += dx
	candidate.Y += dy
	if !g.Board.CanPlace(candidate) {
		return false
	}
	g.Current = candidate
	return true
}

func (g *Game) MoveLeft()  { g.tryMove(-1, 0) }
func (g *Game) MoveRight() { g.tryMove(1, 0) }
func (g *Game) SoftDrop() {
	if g.tryMove(0, 1) {
		g.Score++
	}
}

func (g *Game) Rotate() {
	if g.State != StatePlaying {
		return
	}
	candidate := g.Current.Rotated()
	for _, kick := range []struct{ dx, dy int }{{0, 0}, {-1, 0}, {1, 0}, {-2, 0}, {2, 0}, {0, -1}} {
		candidate.X = g.Current.X + kick.dx
		candidate.Y = g.Current.Y + kick.dy
		if g.Board.CanPlace(candidate) {
			g.Current = candidate
			return
		}
	}
}

func (g *Game) Ghost() piece.Piece {
	ghost := g.Current
	for {
		next := ghost
		next.Y++
		if !g.Board.CanPlace(next) {
			return ghost
		}
		ghost = next
	}
}

func (g *Game) HardDrop() {
	if g.State != StatePlaying {
		return
	}
	target := g.Ghost()
	g.Score += (target.Y - g.Current.Y) * 2
	g.Current = target
	g.lock()
}

func (g *Game) Hold() {
	if g.State != StatePlaying {
		return
	}
	if !g.CanHold {
		return
	}
	g.CanHold = false
	current := g.Current
	if g.Held.Kind == piece.KindEmpty {
		g.Held = current
		g.Current = g.Next
		g.Next = g.newPiece()
		return
	}
	g.Current = g.Held
	g.Current.X = SpawnX
	g.Current.Y = SpawnY
	g.Held = current
}

func (g *Game) lock() {
	for _, c := range g.Current.Cells() {
		if g.Current.Y+c.Y < 0 {
			g.State = StateGameOver
			return
		}
	}
	g.Board.Set(g.Current)
	cleared := g.Board.ClearLines()
	g.addScore(cleared)
	g.Lines += cleared
	g.Level = 1 + g.Lines/LinesPerLevel
	g.FallSpeed = FallSpeedForLevel(g.Level)
	g.spawn()
}

func (g *Game) addScore(cleared int) {
	switch cleared {
	case 1:
		g.Score += 100 * g.Level
	case 2:
		g.Score += 300 * g.Level
	case 3:
		g.Score += 500 * g.Level
	case 4:
		g.Score += 800 * g.Level
	}
}

func FallSpeedForLevel(level int) float64 {
	speed := 1.0 - float64(level-1)*0.08
	if speed < 0.1 {
		speed = 0.1
	}
	return speed
}

// Update advances gravity. dt is seconds. Call once per frame.
func (g *Game) Update(dt float64) {
	if g.State != StatePlaying {
		return
	}
	g.fallTimer += dt
	for g.fallTimer >= g.FallSpeed {
		g.fallTimer -= g.FallSpeed
		if !g.tryMove(0, 1) {
			g.lock()
		}
	}
}

func (g *Game) TogglePause() {
	switch g.State {
	case StatePlaying:
		g.State = StatePaused
	case StatePaused:
		g.State = StatePlaying
	}
}

func (g *Game) Restart() {
	*g = *New(g.generator)
}
