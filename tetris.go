package main

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tetrisTickMsg drives soft gravity while the Tetris pane is active.
type tetrisTickMsg struct{}

func tetrisTickCmd() tea.Cmd {
	return tea.Tick(600*time.Millisecond, func(time.Time) tea.Msg {
		return tetrisTickMsg{}
	})
}

// PieceKind identifies the seven tetrominoes (for rendering / logic).
type PieceKind int

const (
	PieceNone PieceKind = iota
	PieceI
	PieceO
	PieceT
	PieceS
	PieceZ
	PieceJ
	PieceL
)

// Guideline-style 4×4 matrices: spawn state is rotation 0; each step is one
// clockwise rotation (same layouts as tetris.wiki / SRS diagrams).
var pieceDefs = [][][]string{
	// I — horizontal spawn, two flat states + two vertical
	{
		{"....", "####", "....", "...."},
		{"..#.", "..#.", "..#.", "..#."},
		{"....", "....", "####", "...."},
		{".#..", ".#..", ".#..", ".#.."},
	},
	// O
	{
		{".##.", ".##.", "....", "...."},
		{".##.", ".##.", "....", "...."},
		{".##.", ".##.", "....", "...."},
		{".##.", ".##.", "....", "...."},
	},
	// T — flat side down spawn (rounded left in 10-wide well)
	{
		{".#..", "###.", "....", "...."},
		{".#..", ".##.", ".#..", "...."},
		{"....", "###.", ".#..", "...."},
		{".#..", "##..", ".#..", "...."},
	},
	// S
	{
		{".##.", "##..", "....", "...."},
		{".#..", ".##.", "..#.", "...."},
		{"....", ".##.", "##..", "...."},
		{"..#.", ".##.", ".#..", "...."},
	},
	// Z
	{
		{"##..", ".##.", "....", "...."},
		{"..#.", ".##.", ".#..", "...."},
		{"....", "##..", ".##.", "...."},
		{".#..", "##..", ".#..", "...."},
	},
	// J
	{
		{"#...", "###.", "....", "...."},
		{".##.", ".#..", ".#..", "...."},
		{"....", "###.", "..#.", "...."},
		{".#..", ".#..", "##..", "...."},
	},
	// L
	{
		{"..#.", "###.", "....", "...."},
		{".#..", ".#..", ".##.", "...."},
		{"....", "###.", "#...", "...."},
		{"##..", ".#..", ".#..", "...."},
	},
}

func parsePieceGrid(lines []string) (g [4][4]bool) {
	for y := 0; y < 4 && y < len(lines); y++ {
		row := lines[y]
		for x := 0; x < 4 && x < len(row); x++ {
			if row[x] == '#' {
				g[y][x] = true
			}
		}
	}
	return g
}

func pieceShape(kind PieceKind, rot int) [4][4]bool {
	if kind < PieceI || kind > PieceL {
		return [4][4]bool{}
	}
	r := rot % 4
	def := pieceDefs[kind-1][r]
	return parsePieceGrid(def)
}

const (
	standardCols = 10 
	standardRows = 20 
)
func boardSizeFromTerminal(termW, termH int) (cols, rows int) {
	cols = standardCols
	rows = standardRows

	const chromeLines = 9 
	if termH > 0 && termH < chromeLines+standardRows {
		rows = termH - chromeLines
		if rows < 16 {
			rows = 16
		}
	}
	if termW > 0 && termW < 28 {
		_ = termW
	}
	return cols, rows
}

type kick struct{ dx, dy int }

var jlstzKicksCW = [][]kick{
	0: {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	1: {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	2: {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
	3: {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
}

var iKicksCW = [][]kick{
	0: {{0, 0}, {-2, 0}, {1, 0}, {-2, -1}, {1, 2}},
	1: {{0, 0}, {-1, 0}, {2, 0}, {-1, 2}, {2, -1}},
	2: {{0, 0}, {2, 0}, {-1, 0}, {2, 1}, {-1, -2}},
	3: {{0, 0}, {1, 0}, {-2, 0}, {1, -2}, {-2, 1}},
}

type TetrisGame struct {
	Cols, Rows int
	board      [][]PieceKind

	active PieceKind
	next   PieceKind
	rot    int
	px, py int

	gameOver bool
	bag      []PieceKind
	rng      *rand.Rand

	Score int
	Lines int
}

func NewTetrisGame(termW, termH int) *TetrisGame {
	g := &TetrisGame{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	g.Cols, g.Rows = boardSizeFromTerminal(termW, termH)
	g.allocBoard()
	g.refillBag()
	g.active = g.drawFromBag()
	g.next = g.drawFromBag()
	g.resetPieceToTop()
	if g.collidesPiece(g.active, g.px, g.py, g.rot) {
		g.gameOver = true
	}
	return g
}

func (g *TetrisGame) allocBoard() {
	g.board = make([][]PieceKind, g.Rows)
	for y := range g.board {
		g.board[y] = make([]PieceKind, g.Cols)
	}
}

// Resize rebuilds the well when the terminal size changes.
func (g *TetrisGame) Resize(termW, termH int) {
	nc, nr := boardSizeFromTerminal(termW, termH)
	if nc == g.Cols && nr == g.Rows && len(g.board) == g.Rows {
		return
	}
	g.Retry(termW, termH)
}

func (g *TetrisGame) refillBag() {
	g.bag = []PieceKind{PieceI, PieceO, PieceT, PieceS, PieceZ, PieceJ, PieceL}
	g.rng.Shuffle(len(g.bag), func(i, j int) { g.bag[i], g.bag[j] = g.bag[j], g.bag[i] })
}

func (g *TetrisGame) drawFromBag() PieceKind {
	if len(g.bag) == 0 {
		g.refillBag()
	}
	k := g.bag[0]
	g.bag = g.bag[1:]
	return k
}

func (g *TetrisGame) resetPieceToTop() {
	g.rot = 0
	g.px = (g.Cols - 4) / 2
	if g.px < 0 {
		g.px = 0
	}
	g.py = 0
}

func (g *TetrisGame) pieceCells(kind PieceKind, rot, px, py int) [][2]int {
	sh := pieceShape(kind, rot)
	var out [][2]int
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if sh[y][x] {
				out = append(out, [2]int{px + x, py + y})
			}
		}
	}
	return out
}

func (g *TetrisGame) collidesPiece(kind PieceKind, px, py, rot int) bool {
	for _, c := range g.pieceCells(kind, rot, px, py) {
		x, y := c[0], c[1]
		if x < 0 || x >= g.Cols || y < 0 || y >= g.Rows {
			return true
		}
		if g.board[y][x] != PieceNone {
			return true
		}
	}
	return false
}

func (g *TetrisGame) MoveLeft() bool {
	if g.gameOver {
		return false
	}
	if !g.collidesPiece(g.active, g.px-1, g.py, g.rot) {
		g.px--
		return true
	}
	return false
}

func (g *TetrisGame) MoveRight() bool {
	if g.gameOver {
		return false
	}
	if !g.collidesPiece(g.active, g.px+1, g.py, g.rot) {
		g.px++
		return true
	}
	return false
}

func (g *TetrisGame) SoftDrop() bool {
	if g.gameOver {
		return false
	}
	if !g.collidesPiece(g.active, g.px, g.py+1, g.rot) {
		g.py++
		return true
	}
	return false
}

func (g *TetrisGame) RotateCW() bool {
	if g.gameOver {
		return false
	}
	if g.active == PieceO {
		return false
	}
	nextRot := (g.rot + 1) % 4
	var kicks []kick
	if g.active == PieceI {
		kicks = iKicksCW[g.rot]
	} else {
		kicks = jlstzKicksCW[g.rot]
	}
	for _, k := range kicks {
		nx, ny := g.px+k.dx, g.py+k.dy
		if !g.collidesPiece(g.active, nx, ny, nextRot) {
			g.px, g.py = nx, ny
			g.rot = nextRot
			return true
		}
	}
	return false
}

func (g *TetrisGame) TickGravity() {
	if g.gameOver {
		return
	}
	if g.SoftDrop() {
		return
	}
	g.lockPiece()
	g.clearLinesAndScore()
	g.active = g.next
	g.next = g.drawFromBag()
	g.resetPieceToTop()
	if g.collidesPiece(g.active, g.px, g.py, g.rot) {
		g.gameOver = true
	}
}

func (g *TetrisGame) lockPiece() {
	for _, c := range g.pieceCells(g.active, g.rot, g.px, g.py) {
		x, y := c[0], c[1]
		if y >= 0 && y < g.Rows && x >= 0 && x < g.Cols {
			g.board[y][x] = g.active
		}
	}
}

func (g *TetrisGame) clearLinesAndScore() {
	w, h := g.Cols, g.Rows
	kept := make([][]PieceKind, 0, h)
	cleared := 0
	for y := 0; y < h; y++ {
		full := true
		for x := 0; x < w; x++ {
			if g.board[y][x] == PieceNone {
				full = false
				break
			}
		}
		if full {
			cleared++
			continue
		}
		kept = append(kept, g.board[y])
	}
	emptyCount := h - len(kept)
	newBoard := make([][]PieceKind, 0, h)
	for i := 0; i < emptyCount; i++ {
		newBoard = append(newBoard, make([]PieceKind, w))
	}
	newBoard = append(newBoard, kept...)
	g.board = newBoard

	if cleared == 0 {
		return
	}
	g.Lines += cleared
	switch cleared {
	case 1:
		g.Score += 100
	case 2:
		g.Score += 300
	case 3:
		g.Score += 500
	default:
		g.Score += 800 // Tetris (4) or more
	}
}

// Retry clears the well, score, lines, and deals a new bag / active / next.
func (g *TetrisGame) Retry(termW, termH int) {
	g.Cols, g.Rows = boardSizeFromTerminal(termW, termH)
	g.gameOver = false
	g.Score = 0
	g.Lines = 0
	g.allocBoard()
	g.bag = nil
	g.refillBag()
	g.active = g.drawFromBag()
	g.next = g.drawFromBag()
	g.resetPieceToTop()
	if g.collidesPiece(g.active, g.px, g.py, g.rot) {
		g.gameOver = true
	}
}

func (g *TetrisGame) Grid() [][]PieceKind {
	out := make([][]PieceKind, g.Rows)
	for y := 0; y < g.Rows; y++ {
		out[y] = make([]PieceKind, g.Cols)
		copy(out[y], g.board[y])
	}
	if !g.gameOver {
		for _, c := range g.pieceCells(g.active, g.rot, g.px, g.py) {
			x, y := c[0], c[1]
			if y >= 0 && y < g.Rows && x >= 0 && x < g.Cols {
				out[y][x] = g.active
			}
		}
	}
	return out
}

func (g *TetrisGame) RenderLines(cell func(PieceKind) string) []string {
	grid := g.Grid()
	lines := make([]string, len(grid))
	for y := range grid {
		var b strings.Builder
		for x := 0; x < len(grid[y]); x++ {
			b.WriteString(cell(grid[y][x]))
		}
		lines[y] = b.String()
	}
	return lines
}

// NextPreviewLines renders the next piece in a 4×4 mini grid (rotation 0).
func (g *TetrisGame) NextPreviewLines(cell func(PieceKind) string) []string {
	return previewLinesForKind(g.next, cell)
}

func previewLinesForKind(kind PieceKind, cell func(PieceKind) string) []string {
	sh := pieceShape(kind, 0)
	lines := make([]string, 4)
	for y := 0; y < 4; y++ {
		var b strings.Builder
		for x := 0; x < 4; x++ {
			if sh[y][x] {
				b.WriteString(cell(kind))
			} else {
				b.WriteString(cell(PieceNone))
			}
		}
		lines[y] = b.String()
	}
	return lines
}

func (g *TetrisGame) GameOver() bool { return g.gameOver }
