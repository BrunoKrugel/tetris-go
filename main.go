package main

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"
	"tetris-go/internal/board"
	"tetris-go/internal/game"
	"tetris-go/internal/piece"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	cellSize  = 32
	panelW    = 320
	boardPx   = board.Width * cellSize
	winW      = boardPx + panelW
	winH      = board.Height*cellSize + 80
	dasDelay  = 0.25
	dasRepeat = 0.05
)

var (
	colors = map[piece.Kind]color.RGBA{
		piece.KindI: {0, 240, 240, 255},
		piece.KindO: {240, 240, 0, 255},
		piece.KindT: {160, 0, 240, 255},
		piece.KindS: {0, 240, 0, 255},
		piece.KindZ: {240, 0, 0, 255},
		piece.KindJ: {0, 0, 240, 255},
		piece.KindL: {240, 160, 0, 255},
	}
	bgColor    = color.RGBA{20, 20, 30, 255}
	panelBg    = color.RGBA{25, 25, 35, 255}
	emptyColor = color.RGBA{30, 30, 40, 255}
	gridColor  = color.RGBA{45, 45, 60, 255}
	textColor  = color.RGBA{220, 220, 240, 255}
	ghostColor = color.RGBA{255, 255, 255, 50}
	overlayBg  = color.RGBA{0, 0, 0, 180}
	ctrlBg     = color.RGBA{15, 15, 25, 255}
)

var bigFace, smallFace *text.GoTextFace

func init() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}
	bigFace = &text.GoTextFace{Source: src, Size: 28}
	smallFace = &text.GoTextFace{Source: src, Size: 18}
}

type App struct {
	engine                                     *game.Game
	leftHold, rightHold, downHold              bool
	leftTimer, rightTimer, downTimer           float64
	leftRepeat, rightRepeat, downRepeat        bool
	lastPause, lastR, lastC, lastSpace, lastUp bool
}

func NewApp() *App {
	return &App{
		engine: game.New(game.NewSevenBag(uint64(time.Now().UnixNano()))),
	}
}

func (a *App) Layout(w, h int) (int, int) { return winW, winH }

func (a *App) Update() error {
	dt := 1.0 / 60.0
	if a.oneShot(ebiten.KeyP, &a.lastPause) {
		a.engine.TogglePause()
	}
	if a.oneShot(ebiten.KeyR, &a.lastR) {
		a.engine.Restart()
	}
	if a.engine.State != game.StatePlaying {
		return nil
	}
	if a.oneShot(ebiten.KeyC, &a.lastC) {
		a.engine.Hold()
	}
	if a.oneShot(ebiten.KeySpace, &a.lastSpace) {
		a.engine.HardDrop()
	}
	if a.oneShot(ebiten.KeyArrowUp, &a.lastUp) {
		a.engine.Rotate()
	}
	a.das(ebiten.KeyArrowLeft, &a.leftHold, &a.leftTimer, &a.leftRepeat, a.engine.MoveLeft, dt)
	a.das(ebiten.KeyArrowRight, &a.rightHold, &a.rightTimer, &a.rightRepeat, a.engine.MoveRight, dt)
	a.das(ebiten.KeyArrowDown, &a.downHold, &a.downTimer, &a.downRepeat, a.engine.SoftDrop, dt)
	a.engine.Update(dt)
	return nil
}

func (a *App) oneShot(key ebiten.Key, last *bool) bool {
	p := ebiten.IsKeyPressed(key)
	if p && !*last {
		*last = true
		return true
	}
	*last = p
	return false
}

func (a *App) das(key ebiten.Key, hold *bool, timer *float64, repeat *bool, action func(), dt float64) {
	if !ebiten.IsKeyPressed(key) {
		*hold, *repeat = false, false
		*timer = 0
		return
	}
	if !*hold {
		*hold = true
		*timer = 0
		*repeat = false
		action()
		return
	}
	*timer += dt
	if !*repeat {
		if *timer >= dasDelay {
			*repeat = true
			*timer -= dasDelay
			action()
		}
		return
	}
	for *timer >= dasRepeat {
		*timer -= dasRepeat
		action()
	}
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)
	a.drawBoard(screen)
	a.drawPanel(screen)
	a.drawControls(screen)
	switch a.engine.State {
	case game.StatePaused:
		a.drawOverlay(screen, "PAUSADO")
	case game.StateGameOver:
		a.drawOverlay(screen, "FIM DE JOGO\n\npressione R para jogar de novo")
	}
}

func (a *App) drawBoard(screen *ebiten.Image) {
	for y := range board.Height {
		for x := range board.Width {
			c := emptyColor
			if (x+y)%2 == 0 {
				c = gridColor
			}
			drawRect(screen, x*cellSize, y*cellSize, cellSize, cellSize, c)
		}
	}
	for y := range board.Height {
		for x := range board.Width {
			if k := a.engine.Board.Grid[y][x]; k != piece.KindEmpty {
				drawRect(screen, x*cellSize, y*cellSize, cellSize, cellSize, colors[k])
			}
		}
	}
	if a.engine.State == game.StatePlaying {
		ghost := a.engine.Ghost()
		for _, c := range ghost.Cells() {
			px, py := (ghost.X+c.X)*cellSize, (ghost.Y+c.Y)*cellSize
			if py >= 0 {
				vector.StrokeRect(screen, float32(px), float32(py), float32(cellSize), float32(cellSize), 2, ghostColor, false)
			}
		}
		for _, c := range a.engine.Current.Cells() {
			px, py := (a.engine.Current.X+c.X)*cellSize, (a.engine.Current.Y+c.Y)*cellSize
			if py >= 0 {
				drawRect(screen, px, py, cellSize, cellSize, colors[a.engine.Current.Kind])
			}
		}
	}
}

func (a *App) drawPanel(screen *ebiten.Image) {
	drawRect(screen, boardPx, 0, panelW, board.Height*cellSize, panelBg)
	drawText(screen, "PRÓXIMA", boardPx+20, 30, bigFace)
	a.drawPreview(screen, boardPx+20, 55, a.engine.Next)
	drawText(screen, "GUARDADA", boardPx+20, 180, bigFace)
	a.drawPreview(screen, boardPx+20, 205, a.engine.Held)
	drawText(screen, fmt.Sprintf("PONTOS: %d", a.engine.Score), boardPx+20, 340, smallFace)
	drawText(screen, fmt.Sprintf("LINHAS: %d", a.engine.Lines), boardPx+20, 380, smallFace)
	drawText(screen, fmt.Sprintf("NÍVEL: %d", a.engine.Level), boardPx+20, 420, smallFace)
}

func (a *App) drawPreview(screen *ebiten.Image, ox, oy int, p piece.Piece) {
	if p.Kind == piece.KindEmpty {
		return
	}
	cells := p.Cells()
	if len(cells) == 0 {
		return
	}
	minX, minY := cells[0].X, cells[0].Y
	maxX, maxY := cells[0].X, cells[0].Y
	for _, c := range cells[1:] {
		if c.X < minX {
			minX = c.X
		}
		if c.Y < minY {
			minY = c.Y
		}
		if c.X > maxX {
			maxX = c.X
		}
		if c.Y > maxY {
			maxY = c.Y
		}
	}
	cw := cellSize / 2
	for _, c := range cells {
		drawRect(screen, ox+(c.X-minX)*cw, oy+(c.Y-minY)*cw, cw, cw, colors[p.Kind])
	}
}

func (a *App) drawControls(screen *ebiten.Image) {
	cy := board.Height * cellSize
	drawRect(screen, 0, cy, winW, 80, ctrlBg)
	drawText(screen, "setas: mover  ↑: girar  ↓: descer devagar  espaço: soltar", 10, cy+10, smallFace)
	drawText(screen, "c: guardar peça  p: pausar  r: reiniciar", 10, cy+34, smallFace)
}

func (a *App) drawOverlay(screen *ebiten.Image, msg string) {
	drawRect(screen, 0, 0, winW, winH, overlayBg)
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		adv, _ := text.Measure(line, bigFace, 0)
		x := int((float64(winW) - adv) / 2)
		y := winH/2 - len(lines)*20 + i*40
		drawText(screen, line, x, y, bigFace)
	}
}

func main() {
	ebiten.SetWindowSize(winW, winH)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(NewApp()); err != nil {
		panic(err)
	}
}

func drawRect(screen *ebiten.Image, x, y, w, h int, clr color.Color) {
	vector.FillRect(screen, float32(x), float32(y), float32(w), float32(h), clr, false)
}

func drawText(screen *ebiten.Image, txt string, x, y int, f *text.GoTextFace) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(textColor)
	text.Draw(screen, txt, f, op)
}
